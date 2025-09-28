package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hotel-system/src/constants"
	"hotel-system/src/serializers"
	hotelsystem "hotel-system/src/types/hotelsystem"
	"hotel-system/src/utils"
	"hotel-system/src/validators"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

func (s *Service) BookHotel(w http.ResponseWriter, r *http.Request) {
	var bookHotelRequest hotelsystem.BookHotelRequest
	err := json.NewDecoder(r.Body).Decode(&bookHotelRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idempotencyKey := r.Header.Get(constants.IdempotencyKeyHeader)
	resp, err := s.resolveIdempotency(w, &bookHotelRequest, idempotencyKey)
	if resp != nil || err != nil {
		return
	}

	if err = validators.ValidateBookHotelRequest(&bookHotelRequest); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Start transaction: get hotel (lock), update rooms, create booking
	tx, err := s.storageService.BeginTransaction(context.Background())
	defer utils.RollbackOrCommitTransaction(tx, &err)
	hotelForUpdate, err := s.storageService.GetHotelForUpdate(tx, bookHotelRequest.HotelID)
	if err != nil {
		log.Println("Error getting hotel for update:", err)
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	if hotelForUpdate.AvailableRooms < int(bookHotelRequest.NumRooms) {
		sendJsonResponse(w, serializers.BookHotelResponseSerializer("", 0, nil, 0, "Not enough rooms available"))
		return
	}

	//time.Sleep(5 * time.Second) // Simulate some processing delay

	err = s.storageService.DecrementHotelRoomsTx(tx, bookHotelRequest.HotelID, int(bookHotelRequest.NumRooms))
	if err != nil {
		log.Println("Error updating hotel room count:", err)
		http.Error(w, "Failed to update hotel room count", http.StatusInternalServerError)
		return
	}

	userId := r.Context().Value("user_id").(int)
	booking := serializers.BookingSerializer(&bookHotelRequest, userId)
	bookingId, err := s.storageService.CreateBookingTx(tx, booking)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	totalCost := hotelForUpdate.CostPerNight * float32(bookHotelRequest.NumRooms) * float32(bookHotelRequest.NumDays)
	//params := serializers.CreateStripeCheckoutSessionParams(bookingId, hotelForUpdate, totalCost)
	paymentStatusUrl := "http://localhost:6000/status"

	paymentId := uuid.New().String()
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, price_1234) of the product you want to sell
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:          stripe.String(stripe.CurrencyUSD),
					UnitAmountDecimal: stripe.Float64(float64(totalCost * 100)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Hotel Booking for " + hotelForUpdate.Name),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			constants.BookingIdField: fmt.Sprintf("%d", bookingId),
			"payment_id":             paymentId,
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		CustomerEmail: stripe.String("k@gmail.com"),
		SuccessURL:    stripe.String(paymentStatusUrl),
		CancelURL:     stripe.String(paymentStatusUrl),
	}

	checkoutSession, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
		http.Error(w, "Failed to create stripe checkout session", http.StatusInternalServerError)
		return
	}

	newPayment := serializers.CreatePaymentSerializer(paymentId, bookingId, checkoutSession.ID, totalCost)
	err = s.storageService.CreatePayment(newPayment)
	if err != nil {
		log.Println("Error creating payment entry:", err)
		sendJsonResponse(w, hotelsystem.GenericSuccessResponse{
			Status:  "FAIL",
			Message: "Failed to create payment entry",
		})
		return
	}

	CheckoutUrl = checkoutSession.URL // for testing

	resp = serializers.BookHotelResponseSerializer(checkoutSession.URL, bookingId, booking, totalCost, "Booking request created successfully")
	idempotencyEntry := serializers.CreateIdempotencyKeySerializer(&bookHotelRequest, resp, userId, idempotencyKey)
	err = s.storageService.CreateIdempotencyKey(idempotencyEntry)
	if err != nil {
		log.Println("Error creating idempotency key:", err)
	}
	sendJsonResponse(w, &resp)
}

var CheckoutUrl string // for testing

// Redirects the client to the checkout URL set during booking creation.
// Returns 404 if the URL hasn't been set yet.
func (s *Service) ClientCheckoutRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CheckoutUrl:", CheckoutUrl)
	if CheckoutUrl == "" {
		http.Error(w, "Checkout URL not set", http.StatusNotFound)
		return
	}
	// Use 303 See Other to redirect after a non-GET safely
	http.Redirect(w, r, CheckoutUrl, http.StatusSeeOther)
}

func (s *Service) GetBookingDetailsById(w http.ResponseWriter, r *http.Request) {
	var getBookingRequest hotelsystem.GetBookingDetailsRequest
	err := json.NewDecoder(r.Body).Decode(&getBookingRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validators.ValidateGetBookingDetailsByIdRequest(&getBookingRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	booking, err := s.storageService.GetBookingById(getBookingRequest.BookingID)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	payment, err := s.storageService.GetPaymentByBookingId(getBookingRequest.BookingID)
	if err != nil {
		log.Println("Error getting booking details:", err)
	}
	var (
		status    constants.PaymentStatus
		totalCost float32
	)
	if payment != nil {
		status = payment.Status
		totalCost = payment.Amount
	}
	sendJsonResponse(w, hotelsystem.GetBookingDetailsResponse{
		BookingID:     int64(booking.BookingID),
		NumRooms:      int32(booking.NumberOfRooms),
		NumDays:       int32(booking.NumberOfDays),
		CheckInDate:   booking.CheckInDate.Format(constants.DateFormat),
		CheckOutDate:  booking.CheckOutDate.Format(constants.DateFormat),
		TotalCost:     totalCost,
		Status:        string(booking.Status),
		PaymentStatus: string(status),
		BookingTime:   booking.BookingTime.String(),
	})
}

func (s *Service) resolveIdempotency(w http.ResponseWriter, bookHotelRequest *hotelsystem.BookHotelRequest, idempotencyKey string) (*hotelsystem.BookHotelResponse, error) {
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header is required", http.StatusBadRequest)
		return nil, errors.New("missing idempotency key")
	}
	resp, err := s.checkIdempotencyForBookingAndReturnResponse(
		idempotencyKey,
		bookHotelRequest,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	if resp != nil {
		sendJsonResponse(w, resp)
		return resp, nil
	}
	return nil, nil
}

func (s *Service) checkIdempotencyForBookingAndReturnResponse(
	idempotencyKey string,
	bookHotelRequest *hotelsystem.BookHotelRequest,
) (*hotelsystem.BookHotelResponse, error) {
	idempotencyEntry, err := s.storageService.GetIdempotentPayloadByKey(idempotencyKey)
	if err != nil {
		return nil, err
	}
	if idempotencyEntry == nil {
		return nil, nil // key not found
	}
	incomingPayloadBytes, _ := json.Marshal(bookHotelRequest)
	idempotentRequestPayload := idempotencyEntry.RequestPayload

	// Normalize both JSONs
	normalizeJSON := func(b []byte) ([]byte, error) {
		var temp map[string]interface{}
		if err := json.Unmarshal(b, &temp); err != nil {
			return nil, err
		}
		return json.Marshal(temp)
	}

	normStored, err1 := normalizeJSON(idempotentRequestPayload)
	normIncoming, err2 := normalizeJSON(incomingPayloadBytes)

	if err1 != nil || err2 != nil {
		return nil, errors.New("invalid request payload")
	}

	if !bytes.Equal(normStored, normIncoming) {
		return nil, errors.New("same idempotent key but different request")
	}

	var bookHotelResponse hotelsystem.BookHotelResponse
	if err = json.Unmarshal(idempotencyEntry.ResponsePayload, &bookHotelResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stored response: %w", err)
	}
	return &bookHotelResponse, nil

}
