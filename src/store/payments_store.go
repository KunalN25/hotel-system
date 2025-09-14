package store

import (
	"database/sql"
	"hotel-system/src/constants"
	"strings"
)

func (ds *dataStore) CreatePayment(payment *Payment) error {
	query := `
		INSERT INTO public.payments 
		(booking_id, order_id, amount, currency, status, created_at, checkout_session_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING booking_id
	`

	err := ds.db.QueryRow(
		query,
		payment.BookingID,
		payment.OrderID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CreatedAt,
		payment.CheckoutSessionId,
	)

	if err != nil {
		return err.Err()
	}

	return nil
}

func (ds *dataStore) UpdatePaymentStatus(paymentId int, status constants.PaymentStatus) error {
	query := "UPDATE public.payments SET status = $1 WHERE id = $2"
	_, err := ds.db.Exec(query, status, paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetPaymentByCheckoutSessionId(checkoutSessionid string) (*Payment, error) {
	query := "SELECT id, booking_id, order_id, amount, status, checkout_session_id FROM public.payments WHERE order_id = $1"
	var payment Payment
	err := ds.db.QueryRow(query, checkoutSessionid).Scan(&payment.ID, &payment.BookingID, &payment.OrderID, &payment.Amount, &payment.Status, &payment.CheckoutSessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // payment does not exist
		}
		return nil, err // other DB error
	}
	return &payment, nil // payment exists
}

func (ds *dataStore) GetPaymentByBookingId(bookingId int64) (*Payment, error) {
	query := "SELECT " + strings.Join(PaymentsTableColumns, ", ") + " FROM public.payments WHERE booking_id = $1"
	var payment Payment
	err := ds.db.QueryRow(query, bookingId).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // payment does not exist
		}
		return nil, err // other DB error
	}
	return &payment, nil // payment exists
}

func (ds *dataStore) UpdatePaymentStatusTx(tx *sql.Tx, paymentId int, status constants.PaymentStatus) error {
	query := "UPDATE public.payments SET status = $1 WHERE id = $2"
	_, err := tx.Exec(query, status, paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetPaymentByCheckoutSessionIdTx(tx *sql.Tx, checkoutSessionId string) (*Payment, error) {
	query := "SELECT " + strings.Join(PaymentsTableColumns, ", ") + " FROM public.payments WHERE checkout_session_id = $1 FOR UPDATE"
	var payment Payment
	err := tx.QueryRow(query, checkoutSessionId).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
