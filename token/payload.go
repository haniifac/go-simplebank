package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = fmt.Errorf("invalid token")
	ErrExpiredToken = fmt.Errorf("token is expired")
)

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
	// IssuedAt  time.Time `json:"issued_at"`
	// ExpiredAt time.Time `json:"expired_at"`
	// Issuer    string    `json:"issuer,omitempty"`
	// Subject   string    `json:"subject,omitempty"`
	// Audience  []string  `json:"audience,omitempty"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "simplebank",
			Subject:   username,
			Audience:  []string{"simplebank"},
		},
	}

	return payload, nil
}

// func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
// 	return jwt.NewNumericDate(payload.ExpiredAt), nil
// }

// func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
// 	return jwt.NewNumericDate(payload.IssuedAt), nil
// }

// func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
// 	// Optional: Return nil if not applicable
// 	return nil, nil
// }

// func (payload *Payload) GetIssuer() (string, error) {
// 	return payload.Issuer, nil
// }

// func (payload *Payload) GetSubject() (string, error) {
// 	return payload.Subject, nil
// }

// func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
// 	return payload.Audience, nil
// }
