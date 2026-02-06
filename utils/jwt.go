package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID, employeeNIK string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      userID,
		"employee_nik": employeeNIK,
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
