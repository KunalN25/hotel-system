package constants

import "github.com/stripe/stripe-go/v82"

const BookingIdField = "booking_id"
const MaxBodyBytes = int64(65536)
const DateFormat = "2006-01-02"
const IdempotencyKeyHeader = "Idempotency-Key"

type PaymentStatus string

const (
	PAYMENT_SUCCESS PaymentStatus = "success"
	PAYMENT_FAILED  PaymentStatus = "failed"
	PAYMENT_PENDING PaymentStatus = "pending"
	//	add processing?
)

const (
	StripePaymentSucceeded = "payment_intent.succeeded"
	StripePaymentFailed    = "payment_intent.payment_failed"
	StripePaymentCompleted = "checkout.session.completed"
)

var StripePaymentStatusToPaymentStatusMap = map[stripe.EventType]PaymentStatus{
	stripe.EventTypeCheckoutSessionCompleted: PAYMENT_SUCCESS,
	StripePaymentFailed:                      PAYMENT_FAILED,
}
