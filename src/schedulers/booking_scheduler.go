package scheduler

import (
	"context"
	"fmt"
	"hotel-system/src/store"
	"log"

	"github.com/robfig/cron/v3"
)

type BookingScheduler struct {
	storage store.StorageService
}

func NewBookingScheduler(storage store.StorageService) *BookingScheduler {
	return &BookingScheduler{
		storage: storage,
	}
}

func (bs *BookingScheduler) Start() {
	c := cron.New()
	// Runs every 5 minutes
	c.AddFunc("*/5 * * * *", func() {
		err := bs.ReleaseCompletedBookings()
		if err != nil {
			log.Printf("Failed to release expired bookings: %v", err)
		}
	})
	c.AddFunc("*/1 * * * *", func() {
		if err := bs.ExpireStaleBookings(); err != nil {
			log.Printf("Failed to expire stale bookings: %v", err)
		}
	})

	c.Start()
}

func (bs *BookingScheduler) ReleaseCompletedBookings() error {
	bookings, err := bs.storage.GetCompletedBookings()
	if err != nil {
		return err
	}
	err = bs.releaseBookings(bookings, store.BOOKING_COMPLETED)
	if err != nil {
		return fmt.Errorf("failed to release completed bookings: %w", err)
	}
	return nil
}

func (bs *BookingScheduler) ExpireStaleBookings() error {
	bookings, err := bs.storage.GetExpiredBookings()
	if err != nil {
		return err
	}
	err = bs.releaseBookings(bookings, store.BOOKING_EXPIRED)
	if err != nil {
		return fmt.Errorf("failed to release expired bookings: %w", err)
	}
	return nil
}

func (bs *BookingScheduler) releaseBookings(bookings []*store.Booking, newStatus store.BookingStatus) error {

	for _, b := range bookings {
		ctx := context.Background()
		tx, err := bs.storage.BeginTransaction(ctx)
		if err != nil {
			log.Printf("Failed to begin transaction for booking %d: %v\n", b.BookingID, err)
			continue
		}
		hotel, err := bs.storage.GetHotelForUpdate(tx, int64(b.HotelID))
		if err != nil {
			continue
		}
		fmt.Printf("Releasing booking %d for hotel %s\n", b.BookingID, hotel.Name)
		newRooms := hotel.AvailableRooms + b.NumberOfRooms
		if newRooms <= hotel.TotalRooms {
			err = bs.storage.IncrementHotelRoomsTx(tx, int64(hotel.ID), b.NumberOfRooms)
		}
		if err != nil {
			log.Printf("Failed to increment hotel rooms for booking %d: %v", b.BookingID, err)
			continue
		}
		err = bs.storage.UpdateBookingStatusTx(tx, int64(b.BookingID), newStatus)
		if err != nil {
			log.Printf("Failed to update booking status for booking %d: %v", b.BookingID, err)
			continue
		}
	}

	return nil
}
