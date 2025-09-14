package utils

import (
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RollbackOrCommitTransaction(tx *sql.Tx, err *error) {
	if tx == nil {
		return
	}
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if *err != nil {
		_ = tx.Rollback()
	} else {
		*err = tx.Commit()
	}
}

func NewUuid() string {
	return uuid.New().String()
}
