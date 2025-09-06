// Package hash provides password hashing and verification using bcrypt
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher defines methods for hashing and verifying passwords
type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(password, hashedPassword string) bool
}

// BcryptHasher implements PasswordHasher using bcrypt algorithm
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher with the given cost
func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost: cost}
}

// Hash generates a bcrypt hash for the given password
func (h *BcryptHasher) Hash(password string) (string, error) {
	if h.cost == 0 {
		h.cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Check compares a plaintext password with a hashed password
func (h *BcryptHasher) Check(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
