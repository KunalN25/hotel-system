package errorcodes

type DBError string

const (
	ErrUserAddFailed      = "ERR_USER_ADD_FAILED"
	ErrUserNotFoundByID   = "ERR_USER_NOT_FOUND_BY_ID"
	ErrUserNotFoundByName = "ERR_USER_NOT_FOUND_BY_USERNAME"

	ErrGetHotelsFailed      = "ERR_GET_HOTELS_FAILED"
	ErrAddHotelFailed       = "ERR_ADD_HOTEL_FAILED"
	ErrHotelNotFoundByID    = "ERR_HOTEL_NOT_FOUND_BY_ID"
	ErrUpdateHotelRooms     = "ERR_UPDATE_HOTEL_ROOMS_FAILED"
	ErrMetadataNotFound     = "ERR_METADATA_NOT_FOUND"
	ErrBookingCreateFailed  = "ERR_BOOKING_CREATE_FAILED"
	ErrBookingUpdateFailed  = "ERR_BOOKING_UPDATE_STATUS_FAILED"
	ErrBookingNotFoundByID  = "ERR_BOOKING_NOT_FOUND_BY_ID"
	ErrGetCompletedBookings = "ERR_GET_COMPLETED_BOOKINGS_FAILED"

	ErrPaymentCreateFailed  = "ERR_PAYMENT_CREATE_FAILED"
	ErrPaymentNotFoundByOID = "ERR_PAYMENT_NOT_FOUND_BY_ORDER_ID"
	ErrPaymentStatusUpdate  = "ERR_PAYMENT_UPDATE_STATUS_FAILED"
)
