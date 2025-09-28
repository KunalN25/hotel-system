package services

import (
	"database/sql"
	"encoding/json"
	payments2 "hotel-system/src/payments"
	scheduler "hotel-system/src/schedulers"
	"hotel-system/src/serializers"
	"hotel-system/src/store"
	hotelsystem "hotel-system/src/types/hotelsystem"
	"hotel-system/src/validators"
	"log"
	"net/http"
)

type Service struct {
	storageService store.StorageService
	paymentsClient payments2.Client
}

func NewService(db *sql.DB) *Service {
	storageService := store.NewStore(db)
	paymentsClient := payments2.NewPaymentsClient()
	bs := scheduler.NewBookingScheduler(storageService)
	bs.Start()
	return &Service{storageService: storageService, paymentsClient: paymentsClient}
}

func (s *Service) AddHotel(w http.ResponseWriter, r *http.Request) {
	var addHotelRequest hotelsystem.AddHotelRequest
	err := json.NewDecoder(r.Body).Decode(&addHotelRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.storageService.AddHotel(addHotelRequest)
	if err != nil {
		http.Error(w, "Could not add hotel", http.StatusInternalServerError)
		return
	}
	sendSuccessResponse(w, "Hotel added successfully")
}

func (s *Service) GetHotelsList(w http.ResponseWriter, r *http.Request) {
	var getHotelsListReq hotelsystem.GetHotelsListRequest
	err := json.NewDecoder(r.Body).Decode(&getHotelsListReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if getHotelsListReq.Limit == 0 {
		getHotelsListReq.Limit = 10
	}
	if err = validators.ValidateGetHotelsListRequest(&getHotelsListReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hotels, err := s.storageService.GetHotels(&getHotelsListReq)
	if err != nil {
		http.Error(w, "Could not fetch hotels", http.StatusInternalServerError)
		return
	}
	hotelsListResponse := serializers.HotelsResponseSerializer(hotels)
	sendJsonResponse(w, hotelsListResponse)
}

func (s *Service) GetHotelById(w http.ResponseWriter, r *http.Request) {

	var getHotelByIdReq hotelsystem.GetHotelByIdRequest
	err := json.NewDecoder(r.Body).Decode(&getHotelByIdReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hotel, err := s.storageService.GetHotelById(getHotelByIdReq.HotelID)
	if err != nil {
		http.Error(w, "Could not fetch hotel", http.StatusInternalServerError)
		return
	}
	hotelResponse := serializers.HotelByIdResponseSerializer(hotel)
	sendJsonResponse(w, hotelResponse)
}

func (s *Service) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	var paymentStatusRequest hotelsystem.PaymentStatusRequest
	err := json.NewDecoder(r.Body).Decode(&paymentStatusRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Change to payment by booking id or session id
	payment, err := s.storageService.GetPaymentByBookingId(paymentStatusRequest.BookingID)
	if err != nil {
		log.Println("Error getting payment by booking id:", err)
		http.Error(w, "Could not fetch payment details", http.StatusInternalServerError)
		return
	}

	var bookingStatus string
	booking, err := s.storageService.GetBookingById(paymentStatusRequest.BookingID)
	if err != nil {
		log.Println("Error getting booking details:", err)
	}
	if booking != nil {
		bookingStatus = string(booking.Status)
	}

	status := hotelsystem.PaymentStatusResponse{
		OrderID:       paymentStatusRequest.OrderID,
		PaymentStatus: string(payment.Status),
		BookingStatus: bookingStatus,
		Message:       "Payment status retrieved successfully",
	}
	sendJsonResponse(w, status)
}

func sendSuccessResponse(w http.ResponseWriter, message string) {
	successResponse := hotelsystem.GenericSuccessResponse{
		Status:  "S",
		Message: message,
	}
	sendJsonResponse(w, successResponse)
}

func sendJsonResponse(w http.ResponseWriter, jsonData any) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return

	}
}
