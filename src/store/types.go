package store

import (
	"context"
	"database/sql"
	"hotel-system/src/constants"
	"hotel-system/src/types/hotelsystem"
)

//go:generate mockgen -source=types.go -destination=mocks/mock_storage.go -package=mocks
type StorageService interface {
	AddUser(user *User) (int, error)
	GetUserById(id int) (User, error)
	GetUserByUsername(username string) (*User, error)
	GetHotels(getHotelsListReq *hotelsystem.GetHotelsListRequest) ([]Hotel, error)
	AddHotel(addHotelRequest hotelsystem.AddHotelRequest) error
	GetHotelById(id int64) (Hotel, error)
	UpdateHotelRooms(hotelId int64, newRoomCount int) error
	CreateBooking(booking *Booking) (int64, error)
	UpdateBookingStatus(bookingId int64, status BookingStatus) error
	GetBookingById(bookingId int64) (*Booking, error)
	GetCompletedBookings() ([]*Booking, error)
	GetExpiredBookings() ([]*Booking, error)
	CreatePayment(payment *Payment) error
	GetPaymentByCheckoutSessionId(checkoutSessionId string) (*Payment, error)
	GetPaymentByBookingId(bookingId int64) (*Payment, error)
	UpdatePaymentStatus(paymentId string, status constants.PaymentStatus) error
	GetExpiredPayments() ([]*Payment, error)

	UpdateBookingStatusTx(tx *sql.Tx, bookingId int64, status BookingStatus) error
	GetBookingByIdTx(tx *sql.Tx, bookingId int64) (Booking, error)
	CreateBookingTx(tx *sql.Tx, booking *Booking) (int64, error)

	GetHotelForUpdate(tx *sql.Tx, hotelId int64) (Hotel, error)
	UpdateHotelRoomsTx(tx *sql.Tx, hotelId int64, newRoomCount int) error
	DecrementHotelRoomsTx(tx *sql.Tx, hotelId int64, numRooms int) error
	IncrementHotelRoomsTx(tx *sql.Tx, hotelId int64, numRooms int) error
	GetHotelByIdTx(tx *sql.Tx, hotelId int64) (Hotel, error)

	GetPaymentByCheckoutSessionIdTx(tx *sql.Tx, checkoutSessionId string) (*Payment, error)
	UpdatePaymentStatusTx(tx *sql.Tx, paymentId string, status constants.PaymentStatus) error
	GetPaymentByIdTx(tx *sql.Tx, paymentId string) (*Payment, error)

	GetIdempotentPayloadByKey(key string) (*IdempotencyKey, error)
	CreateIdempotencyKey(ik *IdempotencyKey) error

	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}
