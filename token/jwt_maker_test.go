package token

import (
	"testing"
	"time"

	"github.com/dpsigor/cheatsheet-golang-postgres/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestJWTToken(t *testing.T) {
	var tests = []struct {
		name        string
		username    string
		duration    time.Duration
		err         ErrToken
		verifyToken func(t *testing.T, username string, issuedAt time.Time, expiresAt time.Time, duration time.Duration, payload *Payload, err error)
	}{
		{
			name:     "OK",
			duration: time.Second,
			username: util.RandomOwner(),
			verifyToken: func(t *testing.T, username string, issuedAt time.Time, expiresAt time.Time, duration time.Duration, payload *Payload, err error) {
				require.NoError(t, err)
				require.Equal(t, username, payload.Username)
				require.NotEmpty(t, payload.ID)
				require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Millisecond)
				require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Millisecond)
			},
		},
		{
			name:     "Expired",
			duration: -time.Second,
			username: util.RandomOwner(),
			verifyToken: func(t *testing.T, username string, issuedAt time.Time, expiresAt time.Time, duration time.Duration, payload *Payload, err error) {
				require.EqualError(t, err, ErrExpiredToken.Error())
				require.Nil(t, payload)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewJWTMaker(util.RandomString(32))
			require.NoError(t, err)
			token, err := m.CreateToken(tt.username, tt.duration)
			require.NoError(t, err)

			issuedAt := time.Now()
			expiresAt := issuedAt.Add(tt.duration)

			payload, err := m.VerifyToken(token)
			tt.verifyToken(t, tt.username, issuedAt, expiresAt, tt.duration, payload, err)
		})
	}
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := newPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
