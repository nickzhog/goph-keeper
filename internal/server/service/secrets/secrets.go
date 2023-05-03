package secrets

import (
	"crypto/rsa"
	"errors"

	"github.com/nickzhog/goph-keeper/pkg/encryption"
)

var ErrNotFound = errors.New("not found")

const (
	TypeAccount = "ACCOUNT"
	TypeBinary  = "BINARY"
	TypeNote    = "NOTE"
	TypeCard    = "CARD"
)

type AbstractSecret struct {
	ID          string
	UserID      string
	Title       string
	SType       string
	Data        []byte
	IsEncrypted bool
}

func NewSecret(id, userID, title, stype string, data []byte) *AbstractSecret {
	return &AbstractSecret{
		ID:     id,
		UserID: userID,
		Title:  title,
		SType:  stype,
		Data:   data,
	}
}

func (s *AbstractSecret) Encrypt(key *rsa.PublicKey) error {
	if s.IsEncrypted {
		return nil
	}

	encrypted, err := encryption.EncryptData(s.Data, key)
	if err != nil {
		return err
	}

	s.IsEncrypted = true
	s.Data = encrypted
	return nil
}

func (s *AbstractSecret) Decrypt(key *rsa.PrivateKey) error {
	if !s.IsEncrypted {
		return nil
	}

	decrypted, err := encryption.DecryptData(s.Data, key)
	if err != nil {
		return err
	}

	s.Data = decrypted

	return nil
}
