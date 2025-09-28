package store

import (
	"encoding/json"
	"hotel-system/src/constants"
	"time"
)

const SchemaName = "public"

const HotelTableName = "hotel"
const UserTableName = "user"
const BookingTableName = "booking"
const PaymentsTableName = "payments"
const IdempotencyKeyTableName = "idempotency_keys"

type Hotel struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	AvailableRooms int      `json:"available_rooms"`
	TotalRooms     int      `json:"total_rooms"`
	Street         string   `json:"street"`
	Landmark       string   `json:"landmark"`
	Locality       string   `json:"locality"`
	City           string   `json:"city"`
	Pincode        int      `json:"pincode"`
	State          string   `json:"state"`
	ImageUrls      []string `json:"image_urls"`
	CostPerNight   float32  `json:"cost_per_night"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BookingStatus string

const (
	BOOKING_PENDING   BookingStatus = "pending"
	BOOKING_CONFIRMED BookingStatus = "confirmed"
	BOOKING_CANCELLED BookingStatus = "cancelled"
	BOOKING_EXPIRED   BookingStatus = "expired" // user did not pay in time
	BOOKING_FAILED    BookingStatus = "failed"
	BOOKING_COMPLETED BookingStatus = "completed" // user completed the stay
)

type Booking struct {
	BookingID     int           `json:"booking_id"`
	HotelID       int           `json:"hotel_id"`
	UserID        int           `json:"user_id"`
	NumberOfRooms int           `json:"number_of_rooms"`
	NumberOfDays  int           `json:"number_of_days"`
	BookingTime   time.Time     `json:"booking_time"`
	CheckInDate   time.Time     `json:"check_in_date"`
	CheckOutDate  time.Time     `json:"check_out_date"`
	Status        BookingStatus `json:"status"` // booking status
}

type Payment struct {
	ID                string                  `db:"id"`
	BookingID         int                     `db:"booking_id"`
	OrderID           string                  `db:"order_id"` // payment order id
	Amount            float32                 `db:"amount"`
	Currency          string                  `db:"currency"`
	Status            constants.PaymentStatus `db:"status"`
	CreatedAt         time.Time               `db:"created_at"`
	CheckoutSessionId string                  `db:"checkout_session_id"`
}

type IdempotencyKey struct {
	Id              string          `db:"id"`
	Key             string          `db:"idempotency_key"`
	UserID          int             `db:"user_id"`
	Endpoint        string          `db:"endpoint"`
	RequestPayload  json.RawMessage `db:"request_payload"`
	ResponsePayload json.RawMessage `db:"response_payload"`
	CreatedAt       time.Time       `db:"created_at"`
}

var BookingsTableColumns = []string{
	"booking_id", "hotel_id", "user_id", "number_of_rooms", "number_of_days",
	"booking_time", "check_in_date", "check_out_date", "status",
}

var HotelTableColumns = []string{
	"id",
	"name",
	"description",
	"available_rooms",
	"total_rooms",
	"street",
	"landmark",
	"locality",
	"city",
	"pincode",
	"state",
	"image_urls",
	"cost_per_night",
}

var PaymentsTableColumns = []string{
	"id", "booking_id", "order_id", "amount", "currency", "status", "created_at", "checkout_session_id",
}

var IdempotencyKeyColumns = []string{
	"id",
	"idempotency_key",
	"user_id",
	"endpoint",
	"request_payload",
	"response_payload",
	"created_at",
}

func NewHotel(
	id int,
	name string,
	description string,
	availableRooms int,
	totalRooms int,
	street string,
	landmark string,
	locality string,
	city string,
	pincode int,
	state string,
	imageUrls []string,
	costPerNight float32,
) Hotel {
	return Hotel{
		ID:             id,
		Name:           name,
		Description:    description,
		AvailableRooms: availableRooms,
		TotalRooms:     totalRooms,
		Street:         street,
		Landmark:       landmark,
		Locality:       locality,
		City:           city,
		Pincode:        pincode,
		State:          state,
		ImageUrls:      imageUrls,
		CostPerNight:   costPerNight,
	}
}

func NewIdempotencyKey(
	id string,
	key string,
	userID int,
	endpoint string,
	reqPayload, respPayload []byte,
) *IdempotencyKey {
	return &IdempotencyKey{
		Id:              id,
		Key:             key,
		UserID:          userID,
		Endpoint:        endpoint,
		RequestPayload:  json.RawMessage(reqPayload),
		ResponsePayload: json.RawMessage(respPayload),
		CreatedAt:       time.Now(),
	}
}
