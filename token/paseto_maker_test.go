package token

import (
	"testing"
	"time"

	"github.com/dpsigor/cheatsheet-golang-postgres/util"
	"github.com/stretchr/testify/require"
)

func TestPaseto(t *testing.T) {
	m, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	dur := time.Second
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(dur)
	token, err := m.CreateToken(username, dur)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := m.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotEmpty(t, payload.ID)
	require.NotEmpty(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Millisecond)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Millisecond)
}
