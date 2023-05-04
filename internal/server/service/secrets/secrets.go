package secrets

import (
	"crypto/rsa"
	"errors"

	secretaccount "github.com/nickzhog/goph-keeper/internal/server/service/secrets/secrets/account"
	secretbinary "github.com/nickzhog/goph-keeper/internal/server/service/secrets/secrets/binary"
	secretcard "github.com/nickzhog/goph-keeper/internal/server/service/secrets/secrets/card"
	secretnote "github.com/nickzhog/goph-keeper/internal/server/service/secrets/secrets/note"
	"github.com/nickzhog/goph-keeper/pkg/encryption"
)

var ErrNotFound = errors.New("not found")
var ErrWrongType = errors.New("wrong type")
var ErrInvalid = errors.New("invalid secret")

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

func (s *AbstractSecret) IsValid() bool {
	switch s.SType {
	case TypeAccount:
		_, err := secretaccount.New(s.Data)
		return err == nil

	case TypeBinary:
		_, err := secretbinary.New(s.Data)
		return err == nil
	case TypeCard:
		_, err := secretcard.New(s.Data)
		return err == nil
	case TypeNote:
		_, err := secretnote.New(s.Data)
		return err == nil

	default:
		return false
	}
}

func NewSecret(id, userID, title, stype string, data []byte) *AbstractSecret {
	secret := &AbstractSecret{
		ID:     id,
		UserID: userID,
		Title:  title,
		SType:  stype,
		Data:   data,
	}

	return secret
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
