// Package krypto provides crypto methods
package krypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash implements root.Hash.
type Hash struct{}

// Encrypt a salted hash for the input string.
func (c *Hash) Encrypt(s string) string {
	saltedBytes := []byte(s)

	hashedBytes, _ := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	hash := string(hashedBytes)

	return hash
}

// Compare string to generated hash.
func (c *Hash) Compare(hash, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)

	err := bcrypt.CompareHashAndPassword(existing, incoming)
	if err != nil {
		return fmt.Errorf("compare hash and password %w", err)
	}

	return nil
}
