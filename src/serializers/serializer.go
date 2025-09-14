package serializers

import (
	"encoding/json"
	"fmt"
	"hotel-system/src/constants"
	"hotel-system/src/store"
	"hotel-system/src/types/hotelsystem"
	"hotel-system/src/utils"
	"log"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v82"
)

func HotelsResponseSerializer(hotels []store.Hotel) hotelsystem.GetHotelsListResponse {
	var hotelsList []*hotelsystem.HotelData
	for _, hotel := range hotels {
		hotelDataItem := &hotelsystem.HotelData{
			ID:           strconv.FormatInt(int64(hotel.ID), 10),
			Name:         hotel.Name,
			Description:  hotel.Description,
			Images:       hotel.ImageUrls,
			CostPerNight: float32(hotel.CostPerNight),
		}
		hotelsList = append(hotelsList, hotelDataItem)
	}
	return hotelsystem.GetHotelsListResponse{
		HotelsList: hotelsList,
	}
}

func HotelByIdResponseSerializer(hotel store.Hotel) *hotelsystem.GetHotelByIdResponse {

	addressString := fmt.Sprintf("%s, %s, %s, %s, %d, %s", hotel.Street, hotel.Landmark, hotel.Locality, hotel.City, hotel.Pincode, hotel.State)

	return &hotelsystem.GetHotelByIdResponse{
		ID:             int64(hotel.ID),
		Name:           hotel.Name,
		Description:    hotel.Description,
		Images:         hotel.ImageUrls,
		CostPerNight:   int64(hotel.CostPerNight),
		AvailableRooms: int64(hotel.AvailableRooms),
		Address:        addressString,
	}
}

func BookingSerializer(bookHotelRequest *hotelsystem.BookHotelRequest, userId int) *store.Booking {
	checkInDate, err := time.Parse("2006-01-02", bookHotelRequest.CheckInDate)
	if err != nil {
		fmt.Printf("Error parsing check-in date: %v\n", err)
		return nil
	}
	booking := &store.Booking{
		HotelID:       int(bookHotelRequest.HotelID),
		UserID:        userId,
		NumberOfRooms: int(bookHotelRequest.NumRooms),
		NumberOfDays:  int(bookHotelRequest.NumDays),
		BookingTime:   time.Now(),
		CheckInDate:   checkInDate,
		CheckOutDate:  checkInDate.AddDate(0, 0, int(bookHotelRequest.NumDays)),
		Status:        store.BOOKING_PENDING,
	}
	return booking
}

func CreateIdempotencyKeySerializer(
	bookHotelRequest *hotelsystem.BookHotelRequest,
	bookHotelResponse *hotelsystem.BookHotelResponse,
	userId int,
	idempotencyKey string,
) *store.IdempotencyKey {
	bookingReqJson, err := json.Marshal(bookHotelRequest)
	if err != nil {
		log.Println("Error marshalling booking request:", err)
	}

	bookingRespJson, err := json.Marshal(bookHotelResponse)
	if err != nil {
		log.Println("Error marshalling booking response:", err)
	}
	return store.NewIdempotencyKey(
		utils.NewUuid(),
		idempotencyKey,
		userId,
		"/bookHotel",
		bookingReqJson,
		bookingRespJson,
	)
}

func BookHotelResponseSerializer(checkoutUrl string, booking *store.Booking, totalCost float32, message string) *hotelsystem.BookHotelResponse {
	var bookingDetails *hotelsystem.BookHotelResponseBookingDetails
	if booking != nil {
		bookingDetails = &hotelsystem.BookHotelResponseBookingDetails{
			BookingID:   int64(booking.BookingID),
			HotelID:     int64(booking.HotelID),
			NumRooms:    int32(booking.NumberOfRooms),
			NumDays:     int32(booking.NumberOfDays),
			CheckInDate: booking.CheckInDate.String(),
			TotalCost:   totalCost,
		}
	}
	return &hotelsystem.BookHotelResponse{
		Success:        true,
		Message:        message,
		BookingDetails: bookingDetails,
		CheckoutURL:    checkoutUrl,
	}
}

func CreatePaymentSerializer(bookingId int64, checkoutSessionid string, totalCost float32) *store.Payment {
	return &store.Payment{
		BookingID:         int(bookingId),
		Amount:            totalCost,
		Currency:          "INR",
		Status:            constants.PAYMENT_PENDING,
		CreatedAt:         time.Now(),
		CheckoutSessionId: checkoutSessionid,
	}
}

func CreateStripeCheckoutSessionParams(bookingId int64, hotel store.Hotel, totalCost float32) *stripe.CheckoutSessionParams {
	paymentStatusUrl := "http://localhost:6000/status"
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, price_1234) of the product you want to sell
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:          stripe.String(stripe.CurrencyINR),
					UnitAmountDecimal: stripe.Float64(float64(totalCost * 100)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Hotel Booking for " + hotel.Name),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				constants.BookingIdField: fmt.Sprintf("%d", bookingId),
			},
		},
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		CustomerEmail: stripe.String("k@gmail.com"),
		SuccessURL:    stripe.String(paymentStatusUrl),
		CancelURL:     stripe.String(paymentStatusUrl),
	}

	params.AddExpand("payment_intent")
	return params
}
