package token

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO Token maker
type PasetoMaker struct {
	key          paseto.V4SymmetricKey
	symmetricKey []byte
}

// CreateToken creates a new token for a specific username and duration
func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := newPayload(username, duration)
	if err != nil {
		return "", err
	}
	token := paseto.NewToken()
	token.SetExpiration(payload.ExpiresAt)
	token.SetIssuedAt(payload.IssuedAt)
	token.Set("payload", payload)
	return token.V4Encrypt(p.key, nil), nil
}

type pasetoClaim struct {
	Payload Payload `json:"payload"`
}

// VerifyToken checks if the token is valid or not
func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	pasetoToken, err := paseto.NewParser().ParseV4Local(p.key, token, nil)
	if err != nil {
		return nil, err
	}
	claim := pasetoClaim{}
	err = json.Unmarshal(pasetoToken.ClaimsJSON(), &claim)
	if err != nil {
		return nil, err
	}
	return &claim.Payload, nil
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		key:          paseto.NewV4SymmetricKey(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}
