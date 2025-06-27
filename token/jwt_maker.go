package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretKeySize = 32
)

var (
	ErrInvalidKeySize = fmt.Errorf("invalid key size: must be at least %d bytes", minSecretKeySize)
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInvalidKeySize
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", &Payload{}, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", &Payload{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
