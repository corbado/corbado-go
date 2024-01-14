package entities

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.RegisteredClaims
	Name        string `json:"name,omitempty"`
	Orig        string `json:"orig,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Version     int    `json:"version,omitempty"`
}
