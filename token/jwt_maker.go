package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// CreateToken creates a new token for a specific username and duration
func (j *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	var token string
	var err error
	payload, err := newPayload(username, duration)
	if err != nil {
		return token, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretKey))
}

// VerifyToken checks if the token is valid or not
func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	jwtToken, err := jwt.ParseWithClaims(token, &payload, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &payload, nil
}

// NewJWTMaker returns an instance of a JWT token maker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}
