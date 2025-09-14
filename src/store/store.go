package store

import (
	"context"
	"database/sql"
	"hotel-system/src/types/hotelsystem"
	"strings"

	"github.com/lib/pq"
)

type dataStore struct {
	db *sql.DB
}

func (ds *dataStore) GetHotels(getHotelsListRequest *hotelsystem.GetHotelsListRequest) ([]Hotel, error) {
	var hotels []Hotel
	query := `
		SELECT id, name, description, available_rooms, total_rooms, street, landmark, 
			   locality, city, pincode, state, image_urls, cost_per_night
		FROM public.hotel
		WHERE name      ILIKE $1 OR
			  street    ILIKE $1 OR
			  landmark  ILIKE $1 OR
			  locality  ILIKE $1 OR
			  city      ILIKE $1 OR
			  state     ILIKE $1
		LIMIT $2 OFFSET $3;
	`

	likePattern := "%" + getHotelsListRequest.SearchQuery + "%"
	rows, err := ds.db.Query(query, likePattern, getHotelsListRequest.Limit, getHotelsListRequest.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id, availableRooms, totalRooms, pincode int
		var costPerNight float32
		var imageUrls []string
		var name, description, street, landmark, locality, city, state string
		err = rows.Scan(&id, &name, &description, &availableRooms, &totalRooms, &street, &landmark, &locality, &city, &pincode, &state, pq.Array(&imageUrls), &costPerNight)
		if err != nil {
			return nil, err
		}
		hotel := Hotel{
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
			ImageUrls:      imageUrls, // Convert string to []string
			CostPerNight:   costPerNight,
		}
		hotels = append(hotels, hotel)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (ds *dataStore) AddHotel(addHotelRequest hotelsystem.AddHotelRequest) error {
	query := "INSERT INTO public.hotel (name, description, available_rooms, total_rooms, street, landmark, locality, city, pincode, state, image_urls, cost_per_night) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	_, err := ds.db.Exec(query, addHotelRequest.Name, addHotelRequest.Description, addHotelRequest.TotalRooms, addHotelRequest.TotalRooms, addHotelRequest.Street, addHotelRequest.Landmark, addHotelRequest.Locality, addHotelRequest.City, addHotelRequest.Pincode, addHotelRequest.State, nil, addHotelRequest.CostPerNight)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetHotelById(hotelId int64) (Hotel, error) {
	var id, availableRooms, totalRooms, pincode int
	var costPerNight float32
	var imageUrls []string
	var name, description, street, landmark, locality, city, state string

	query := "SELECT " + strings.Join(HotelTableColumns, ", ") + " FROM public.hotel WHERE id = $1"
	err := ds.db.QueryRow(query, hotelId).Scan(&id, &name, &description, &availableRooms, &totalRooms, &street, &landmark, &locality, &city, &pincode, &state, pq.Array(&imageUrls), &costPerNight)
	if err != nil {
		return Hotel{}, err
	}
	hotel := Hotel{
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
		ImageUrls:      imageUrls, // Convert string to []string
		CostPerNight:   costPerNight,
	}
	return hotel, nil
}

func (ds *dataStore) UpdateHotelRooms(hotelId int64, newRoomCount int) error {
	query := "UPDATE public.hotel SET available_rooms = $1 WHERE id = $2"
	_, err := ds.db.Exec(query, newRoomCount, hotelId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) CreateBooking(booking *Booking) (int64, error) {
	query := `
		INSERT INTO public.booking 
		(hotel_id, user_id, number_of_rooms, number_of_days, booking_time, check_in_date, check_out_date, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING booking_id
	`

	var bookingID int64
	err := ds.db.QueryRow(
		query,
		booking.HotelID,
		booking.UserID,
		booking.NumberOfRooms,
		booking.NumberOfDays,
		booking.BookingTime,
		booking.CheckInDate,
		booking.CheckOutDate,
		booking.Status,
	).Scan(&bookingID)

	if err != nil {
		return 0, err
	}

	return bookingID, nil
}

func (ds *dataStore) GetUserById(id int) (User, error) {
	query := "SELECT id, username FROM public.user WHERE id = $1"
	var user User
	err := ds.db.QueryRow(query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (ds *dataStore) UpdateBookingStatus(bookingId int64, status BookingStatus) error {
	query := "UPDATE public.booking SET status = $1 WHERE booking_id = $2"
	_, err := ds.db.Exec(query, status, bookingId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetBookingById(bookingId int64) (Booking, error) {
	query := "SELECT booking_id, hotel_id, user_id, number_of_rooms, number_of_days, booking_time, check_in_date, check_out_date, status FROM public.booking WHERE booking_id = $1"
	var b Booking
	err := ds.db.QueryRow(query, bookingId).Scan(
		&b.BookingID,
		&b.HotelID,
		&b.UserID,
		&b.NumberOfRooms,
		&b.NumberOfDays,
		&b.BookingTime,
		&b.CheckInDate,
		&b.CheckOutDate,
		&b.Status,
	)
	if err != nil {
		return Booking{}, err
	}
	return b, nil
}

func (ds *dataStore) GetCompletedBookings() ([]*Booking, error) {

	query := "SELECT " + strings.Join(BookingsTableColumns, ", ") +
		" FROM public.booking WHERE status = 'confirmed' AND check_out_date < NOW()"

	rows, err := ds.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(
			&b.BookingID,
			&b.HotelID,
			&b.UserID,
			&b.NumberOfRooms,
			&b.NumberOfDays,
			&b.BookingTime,
			&b.CheckInDate,
			&b.CheckOutDate,
			&b.Status,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (ds *dataStore) AddUser(user *User) (int, error) {
	query := "INSERT INTO public.user (username, password) VALUES ($1, $2) RETURNING id"
	var userID int
	err := ds.db.QueryRow(query, user.Username, user.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ds *dataStore) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, password FROM public.user WHERE username = $1"
	var user User
	err := ds.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // username does not exist
		}
		return nil, err // other DB error
	}
	return &user, nil // username exists
}

func (ds *dataStore) GetExpiredBookings() ([]*Booking, error) {
	query := "SELECT " + strings.Join(BookingsTableColumns, ", ") +
		" FROM public.booking WHERE status = 'pending' AND booking_time < NOW() - INTERVAL '15 minutes'"

	rows, err := ds.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*Booking
	for rows.Next() {
		var b Booking
		if err = rows.Scan(
			&b.BookingID,
			&b.HotelID,
			&b.UserID,
			&b.NumberOfRooms,
			&b.NumberOfDays,
			&b.BookingTime,
			&b.CheckInDate,
			&b.CheckOutDate,
			&b.Status,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (ds *dataStore) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := ds.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (ds *dataStore) UpdateBookingStatusTx(tx *sql.Tx, bookingId int64, status BookingStatus) error {
	query := "UPDATE public.booking SET status = $1 WHERE booking_id = $2"
	_, err := tx.Exec(query, status, bookingId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetHotelByIdTx(tx *sql.Tx, hotelId int64) (Hotel, error) {
	query := "SELECT " + strings.Join(HotelTableColumns, ", ") + " FROM public.hotel WHERE id = $1"
	var id, availableRooms, totalRooms, pincode int
	var costPerNight float32
	var imageUrls []string
	var name, description, street, landmark, locality, city, state string
	err := tx.QueryRow(query, hotelId).Scan(&id, &name, &description, &availableRooms, &totalRooms, &street, &landmark, &locality, &city, &pincode, &state, pq.Array(&imageUrls), &costPerNight)
	if err != nil {
		return Hotel{}, err
	}
	hotel := NewHotel(
		id,
		name,
		description,
		availableRooms,
		totalRooms,
		street,
		landmark,
		locality,
		city,
		pincode,
		state,
		imageUrls,
		costPerNight,
	)

	return hotel, nil
}

func (ds *dataStore) GetHotelForUpdate(tx *sql.Tx, hotelId int64) (Hotel, error) {
	query := "SELECT " + strings.Join(HotelTableColumns, ", ") + " FROM public.hotel WHERE id = $1 FOR UPDATE"
	var id, availableRooms, totalRooms, pincode int
	var costPerNight float32
	var imageUrls []string
	var name, description, street, landmark, locality, city, state string
	err := tx.QueryRow(query, hotelId).Scan(&id, &name, &description, &availableRooms, &totalRooms, &street, &landmark, &locality, &city, &pincode, &state, pq.Array(&imageUrls), &costPerNight)
	if err != nil {
		return Hotel{}, err
	}
	hotel := NewHotel(
		id,
		name,
		description,
		availableRooms,
		totalRooms,
		street,
		landmark,
		locality,
		city,
		pincode,
		state,
		imageUrls,
		costPerNight,
	)
	return hotel, nil
}

func (ds *dataStore) UpdateHotelRoomsTx(tx *sql.Tx, hotelId int64, newRoomCount int) error {
	query := "UPDATE public.hotel SET available_rooms = $1 WHERE id = $2"
	_, err := tx.Exec(query, newRoomCount, hotelId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) DecrementHotelRoomsTx(tx *sql.Tx, hotelId int64, numRooms int) error {
	query := "UPDATE public.hotel SET available_rooms = available_rooms - $1 WHERE id = $2"
	_, err := tx.Exec(query, numRooms, hotelId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) IncrementHotelRoomsTx(tx *sql.Tx, hotelId int64, numRooms int) error {
	query := "UPDATE public.hotel SET available_rooms = available_rooms + $1 WHERE id = $2"
	_, err := tx.Exec(query, numRooms, hotelId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) CreateBookingTx(tx *sql.Tx, booking *Booking) (int64, error) {
	query := `
		INSERT INTO public.booking 
		(hotel_id, user_id, number_of_rooms, number_of_days, booking_time, check_in_date, check_out_date, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING booking_id
	`

	var bookingID int64
	err := tx.QueryRow(
		query,
		booking.HotelID,
		booking.UserID,
		booking.NumberOfRooms,
		booking.NumberOfDays,
		booking.BookingTime,
		booking.CheckInDate,
		booking.CheckOutDate,
		booking.Status,
	).Scan(&bookingID)

	if err != nil {
		return 0, err
	}

	return bookingID, nil
}

func (ds *dataStore) GetBookingByIdTx(tx *sql.Tx, bookingId int64) (Booking, error) {
	query := "SELECT booking_id, hotel_id, user_id, number_of_rooms, number_of_days, booking_time, check_in_date, check_out_date, status FROM public.booking WHERE booking_id = $1 FOR UPDATE"
	var b Booking
	err := tx.QueryRow(query, bookingId).Scan(
		&b.BookingID,
		&b.HotelID,
		&b.UserID,
		&b.NumberOfRooms,
		&b.NumberOfDays,
		&b.BookingTime,
		&b.CheckInDate,
		&b.CheckOutDate,
		&b.Status,
	)
	if err != nil {
		return Booking{}, err
	}
	return b, nil
}

func NewStore(db *sql.DB) StorageService {
	return &dataStore{db: db}
}
