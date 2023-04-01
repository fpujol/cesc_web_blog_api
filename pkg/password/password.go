package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(email string, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+email), bcrypt.DefaultCost)
	if err!=nil {
		return "",fmt.Errorf("failed to hash password: %w",err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(email string, password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password+email))
}