package services

import (
	"context"
	"encoding/json"
	"fmt"
	"hotel-system/src/constants"
	"hotel-system/src/errorcodes"
	"hotel-system/src/store"
	"hotel-system/src/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

// if anything inside here fails, then refund
func (s *Service) PaymentWebhookHandler(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, constants.MaxBodyBytes)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// Use the secret provided by Stripe CLI for local testing or your webhook endpoint's secret.
	webhookSigningSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSigningSecret == "" {
		log.Println("STRIPE_WEBHOOK_SECRET not set in environment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), webhookSigningSecret)

	if err != nil {
		fmt.Printf("Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	var checkoutSession stripe.CheckoutSession
	err = json.Unmarshal(event.Data.Raw, &checkoutSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Webhook received for event type:", event.Type)
	fmt.Println("Checkout session metadata: ", checkoutSession.ID)

	metadataMap := checkoutSession.Metadata
	if metadataMap == nil || metadataMap[constants.BookingIdField] == "" {
		log.Println("Metadata is missing or booking_id is not present")
		http.Error(w, errorcodes.ErrMetadataNotFound, http.StatusBadRequest)
		return
	}

	bookingId, err := strconv.Atoi(metadataMap[constants.BookingIdField])
	if err != nil {
		log.Println("Error converting booking_id to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paymentId, ok := metadataMap["payment_id"]
	if !ok {
		log.Println("payment_id not found in metadata")
		http.Error(w, errorcodes.ErrMetadataNotFound, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tx, err := s.storageService.BeginTransaction(ctx)
	if err != nil {
		log.Println("Failed to start transaction")
		return
	}
	defer utils.RollbackOrCommitTransaction(tx, &err)

	booking, err := s.storageService.GetBookingByIdTx(tx, int64(bookingId))
	if err != nil {
		log.Println("Error getting booking entry:", err)
		http.Error(w, errorcodes.ErrBookingNotFoundByID, http.StatusNotFound)
		return
	}

	if booking.Status == store.BOOKING_FAILED || booking.Status == store.BOOKING_EXPIRED {
		log.Println("Booking already expired or failed")
		w.WriteHeader(http.StatusOK)
		return
	}

	if booking.Status == store.BOOKING_CONFIRMED {
		log.Println("Booking already confirmed, no further action required")
		w.WriteHeader(http.StatusOK)
		return
	}

	payment, err := s.storageService.GetPaymentByIdTx(tx, paymentId)
	if err != nil {
		log.Println("Error getting payment entry:", err)
		http.Error(w, errorcodes.ErrPaymentNotFoundByOID, http.StatusNotFound)
		return
	}

	if payment.Status == constants.PAYMENT_SUCCESS || payment.Status == constants.PAYMENT_FAILED {
		log.Println("Payment already processed, no further action required")
		w.WriteHeader(http.StatusOK)
		return
	}

	//paymentStatus, _ := constants.StripePaymentStatusToPaymentStatusMap[string(event.Type)]
	paymentStatus := event.Type

	switch paymentStatus {
	case stripe.EventTypeCheckoutSessionCompleted, stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded:
		// Update booking status to confirmed
		booking.Status = store.BOOKING_CONFIRMED
		err = s.storageService.UpdateBookingStatusTx(tx, int64(bookingId), booking.Status)
		if err != nil {
			log.Println("Error updating booking status:", err)
			http.Error(w, "Failed to update booking", http.StatusInternalServerError)
			return
		}

		err = s.storageService.UpdatePaymentStatusTx(tx, payment.ID, constants.PAYMENT_SUCCESS)
		if err != nil {
			log.Println("Error updating payment status:", err)
			http.Error(w, "Failed to update payment status", http.StatusInternalServerError)
			return
		}
	case stripe.EventTypeCheckoutSessionAsyncPaymentFailed:
		var hotel store.Hotel
		hotel, err = s.storageService.GetHotelForUpdate(tx, int64(booking.HotelID))
		if err != nil {
			log.Println("Error getting hotel entry:", err)
			http.Error(w, "Hotel not found", http.StatusNotFound)
			return
		}

		err = s.storageService.IncrementHotelRoomsTx(tx, int64(hotel.ID), booking.NumberOfRooms)
		if err != nil {
			log.Println("Error updating hotel rooms:", err)
			http.Error(w, errorcodes.ErrUpdateHotelRooms, http.StatusInternalServerError)
			return
		}

		err = s.storageService.UpdateBookingStatusTx(tx, int64(bookingId), store.BOOKING_FAILED)
		if err != nil {
			log.Println("Error updating booking status:", err)
			http.Error(w, errorcodes.ErrBookingUpdateFailed, http.StatusInternalServerError)
			return
		}

		err = s.storageService.UpdatePaymentStatusTx(tx, payment.ID, constants.PAYMENT_FAILED)
		if err != nil {
			log.Println("Error updating payment status:", err)
			http.Error(w, errorcodes.ErrPaymentStatusUpdate, http.StatusInternalServerError)
			return
		}
	}
}

/*

Above webhook workflow

- Receive HTTP POST webhook and JSON-decode payload; if decode fails, return 400.
- Verify request signature; if invalid, return 400/401.

mark payment as  processing

-
- Extract booking_id from payload metadata; if missing/invalid, return 400.
- Start a database transaction for locking and atomic updates.
- Fetch booking by booking_id; if not found, return 404.
- If booking already CONFIRMED or FAILED/EXPIRED, return 200 (idempotent).
- Look up payment by checkout session ID from payload; if missing/not found, return 404.
- if payment already failed or success then return 200 (idempotent)
- If payment status is SUCCESS:
    - Set booking to CONFIRMED and payment to SUCCESS.

- If payment status is FAILED:
    - Lock hotel, return reserved rooms to inventory, set booking to FAILED, set payment to FAILED.

- Commit on success; on any failure, rollback and return appropriate 4xx/5xx.
- For any other payment status, make no state changes.






Refunds query

CREATE TABLE refunds (
    id SERIAL PRIMARY KEY,

    payment_id VARCHAR NOT NULL,        -- Razorpay payment ID
    booking_id INTEGER NOT NULL,        -- Your internal booking ID

    refund_id VARCHAR,                  -- Razorpay refund ID (after successful API call)
    refund_status VARCHAR NOT NULL,     -- INITIATED / PROCESSED / FAILED
    refund_reason TEXT,                 -- Optional reason (e.g., overbooking, system error)

    amount INTEGER NOT NULL,            -- Amount to refund (in paise)

    attempts INTEGER DEFAULT 0,         -- Number of refund attempts
    last_attempt_at TIMESTAMP,          -- Last refund attempt time

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_refund_payment FOREIGN KEY (payment_id) REFERENCES payments(payment_id),
    CONSTRAINT fk_refund_booking FOREIGN KEY (booking_id) REFERENCES bookings(booking_id)
);

*/
