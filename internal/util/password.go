package util

import "golang.org/x/crypto/bcrypt"

var (
	passwordIterations = 12
	passwordMinLength  = 8
	passwordMaxLength  = 1024
)

// HashPassword creates a bcrypt hash of the given password.
func HashPassword(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, passwordIterations)
	return bytes, err
}

// CheckPasswordHash checks if the given password matches the hash.
func CheckPasswordHash(password []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
