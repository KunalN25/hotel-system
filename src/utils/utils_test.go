package utils

import (
	"testing"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	originalPassword := "MySecurePass123!"

	// Hash the password
	hashedPassword, err := HashPassword(originalPassword)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	// Check the password with the hash
	if !CheckPasswordHash(originalPassword, hashedPassword) {
		t.Error("Password check failed: expected match")
	}

	// Check with an incorrect password
	wrongPassword := "WrongPass456!"
	if CheckPasswordHash(wrongPassword, hashedPassword) {
		t.Error("Password check passed unexpectedly with wrong password")
	}
}
