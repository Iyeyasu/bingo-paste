package auth

import "golang.org/x/crypto/bcrypt"

var (
	// PasswordIterations determins how many bcrypt iterations to use.
	PasswordIterations = 12
)

// HashPassword creates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordIterations)
	return string(bytes), err
}

// CheckPasswordHash checks if the given password matches the hash.
func CheckPasswordHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
