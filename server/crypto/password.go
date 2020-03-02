package crypto

import (
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts a plaintext string and returns the hashed version in base64
func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(hashed)
}
