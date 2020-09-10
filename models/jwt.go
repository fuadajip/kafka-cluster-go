package models

import (
	"github.com/dgrijalva/jwt-go"
)

// UserJWT is struct to map jwt from user-service and used as transactional token for authentication
type UserJWT struct {
	UUID        *string `json:"uuid"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	jwt.StandardClaims
}
