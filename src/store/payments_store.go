package store

import (
	"database/sql"
	"fmt"
	"hotel-system/src/constants"
	"strings"
)

// GetExpiredPayments returns all pending payments created more than 15 minutes ago
func (ds *dataStore) GetExpiredPayments() ([]*Payment, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s.%s WHERE status = $1 AND created_at < NOW() - INTERVAL '15 minutes'",
		strings.Join(PaymentsTableColumns, ", "),
		SchemaName,
		PaymentsTableName,
	)
	rows, err := ds.db.Query(query, constants.PAYMENT_PENDING)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.ID,
			&payment.BookingID,
			&payment.OrderID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.CreatedAt,
			&payment.CheckoutSessionId,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (ds *dataStore) CreatePayment(payment *Payment) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (id, booking_id, order_id, amount, currency, status, created_at, checkout_session_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING booking_id", SchemaName, PaymentsTableName)
	err := ds.db.QueryRow(
		query,
		payment.ID,
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

func (ds *dataStore) UpdatePaymentStatus(paymentId string, status constants.PaymentStatus) error {
	query := fmt.Sprintf("UPDATE %s.%s SET status = $1 WHERE id = $2", SchemaName, PaymentsTableName)
	_, err := ds.db.Exec(query, status, paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetPaymentByCheckoutSessionId(checkoutSessionid string) (*Payment, error) {
	query := fmt.Sprintf("SELECT id, booking_id, order_id, amount, status, checkout_session_id FROM %s.%s WHERE order_id = $1", SchemaName, PaymentsTableName)
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
	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE booking_id = $1", strings.Join(PaymentsTableColumns, ", "), SchemaName, PaymentsTableName)
	var payment Payment
	err := ds.db.QueryRow(query, bookingId).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.CheckoutSessionId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // payment does not exist
		}
		return nil, err // other DB error
	}
	return &payment, nil // payment exists
}

func (ds *dataStore) UpdatePaymentStatusTx(tx *sql.Tx, paymentId string, status constants.PaymentStatus) error {
	query := fmt.Sprintf("UPDATE %s.%s SET status = $1 WHERE id = $2", SchemaName, PaymentsTableName)
	_, err := tx.Exec(query, status, paymentId)
	if err != nil {
		return err
	}
	return nil
}

func (ds *dataStore) GetPaymentByCheckoutSessionIdTx(tx *sql.Tx, checkoutSessionId string) (*Payment, error) {
	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE checkout_session_id = $1 FOR UPDATE", strings.Join(PaymentsTableColumns, ", "), SchemaName, PaymentsTableName)
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

// GetPaymentByIdTx fetches a payment by id using the provided transaction
func (ds *dataStore) GetPaymentByIdTx(tx *sql.Tx, paymentId string) (*Payment, error) {
	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE id = $1 FOR UPDATE", strings.Join(PaymentsTableColumns, ", "), SchemaName, PaymentsTableName)
	var payment Payment
	err := tx.QueryRow(query, paymentId).Scan(
		&payment.ID,
		&payment.BookingID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.CheckoutSessionId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}
