package util

import "golang.org/x/crypto/bcrypt"

var (
	passwordIterations = 12
	passwordMinLength  = 8
	passwordMaxLength  = 1024
)

// HashPassword creates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordIterations)
	return string(bytes), err
}

// CheckPasswordHash checks if the given password matches the hash.
func CheckPasswordHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
