package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func (ds *dataStore) CreateIdempotencyKey(ik *IdempotencyKey) error {
	query := fmt.Sprintf("INSERT INTO %s.%s (id, idempotency_key, user_id, endpoint, request_payload, response_payload, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", SchemaName, IdempotencyKeyTableName)
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
	query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE idempotency_key = $1", strings.Join(IdempotencyKeyColumns, ", "), SchemaName, IdempotencyKeyTableName)
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
