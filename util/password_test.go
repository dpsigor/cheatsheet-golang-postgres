package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pw := RandomString(6)

	hs1, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hs1)

	err = CheckPassword(pw, hs1)
	require.NoError(t, err)

	wrongePw := pw + "x"
	err = CheckPassword(wrongePw, hs1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// The same pw creates different hashes, as the salt is randomized
	hs2, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hs2)

	require.NotEqual(t, hs1, hs2)
}
