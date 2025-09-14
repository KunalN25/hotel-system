package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (ds *dataStore) CreateIdempotencyKey(ik *IdempotencyKey) error {
	query := `
	INSERT INTO idempotency_keys (id, idempotency_key, user_id, endpoint, request_payload, response_payload, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	_, err := ds.db.Exec(
		query,
		ik.Id,
		ik.Key,
		ik.UserID,
		ik.Endpoint,
		ik.RequestPayload,
		ik.ResponsePayload,
		ik.CreatedAt,
	)

	return err
}

func (ds *dataStore) GetIdempotentPayloadByKey(key string) (*IdempotencyKey, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM idempotency_keys
		WHERE idempotency_key = $1
	`, strings.Join(IdempotencyKeyColumns, ", "))

	row := ds.db.QueryRow(query, key)

	var ik IdempotencyKey
	err := row.Scan(
		&ik.Id,
		&ik.Key,
		&ik.UserID,
		&ik.Endpoint,
		&ik.RequestPayload,
		&ik.ResponsePayload,
		&ik.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // return nil if no rows found
		}
		return nil, err
	}

	return &ik, nil
}
