package validators

import (
	"errors"
	"fmt"
	pb2 "hotel-system/src/pb"
	"hotel-system/src/types/hotelsystem"
	"time"
)

func validateUserInfo(username, password string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

func ValidateLoginRequest(req *pb2.LoginRequest) error {
	return validateUserInfo(req.Username, req.Password)
}

func ValidateRegisterRequest(req *pb2.RegisterRequest) error {
	return validateUserInfo(req.Username, req.Password)
}

func ValidateGetHotelsListRequest(req *hotelsystem.GetHotelsListRequest) error {
	if req.Limit < 0 {
		return errors.New("limit must be >= 0")
	}
	if req.Offset < 0 {
		return errors.New("offset must be >= 0")
	}
	if req.SearchQuery == "" {
		return errors.New("search_query cannot be empty")
	}
	// (Optionally) enforce a maximum for limit, e.g.:
	// if req.Limit > 100 {
	//     return fmt.Errorf("limit cannot exceed 100, got %d", req.Limit)
	// }
	return nil
}

// ValidateBookHotelRequest checks that hotel_id, user_id (if still used),
// num_rooms, and num_days are all positive integers.
// Once you remove user_id (after adding auth), drop that check.
func ValidateBookHotelRequest(req *hotelsystem.BookHotelRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.NumRooms <= 0 {
		return fmt.Errorf("num_rooms must be > 0, got %d", req.NumRooms)
	}
	if req.NumDays <= 0 {
		return fmt.Errorf("num_days must be > 0, got %d", req.NumDays)
	}
	if req.CheckInDate == "" {
		return errors.New("check_in_date cannot be empty")
	}
	checkInDateTime, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return fmt.Errorf("check_in_date must be in YYYY-MM-DD format, got %s", req.CheckInDate)
	}
	today := time.Now().Truncate(24 * time.Hour)
	if checkInDateTime.Before(today) {
		return fmt.Errorf("check_in_date cannot be in the past")
	}
	return nil
}

func ValidateGetBookingDetailsByIdRequest(req *hotelsystem.GetBookingDetailsRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}
	if req.BookingID <= 0 {
		return fmt.Errorf("booking_id must be > 0, got %d", req.BookingID)
	}
	return nil
}
