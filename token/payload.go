package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrToken is a token validation error
type ErrToken error

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken ErrToken = errors.New("token is invalid")
	ErrExpiredToken ErrToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// newPayload creates a new token payload with a specific username and duration
func newPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	now := time.Now()

	if now.After(payload.ExpiresAt) {
		return ErrExpiredToken
	}

	if now.Before(payload.IssuedAt) {
		return ErrInvalidToken
	}

	return nil
}
