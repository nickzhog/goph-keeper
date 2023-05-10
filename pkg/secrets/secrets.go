package secrets

import (
	"crypto/rsa"
	"errors"

	"github.com/nickzhog/goph-keeper/pkg/encryption"
	secretaccount "github.com/nickzhog/goph-keeper/pkg/secrets/account"
	secretbinary "github.com/nickzhog/goph-keeper/pkg/secrets/binary"
	secretcard "github.com/nickzhog/goph-keeper/pkg/secrets/card"
	secretnote "github.com/nickzhog/goph-keeper/pkg/secrets/note"
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

func (s *AbstractSecret) Validate() error {
	switch s.SType {
	case TypeAccount:
		_, err := secretaccount.Unmarshal(s.Data)
		if err != nil {
			return ErrInvalid
		}
		return nil

	case TypeBinary:
		_, err := secretbinary.Unmarshal(s.Data)
		if err != nil {
			return ErrInvalid
		}
		return nil
	case TypeCard:
		_, err := secretcard.Unmarshal(s.Data)
		if err != nil {
			return ErrInvalid
		}
		return nil
	case TypeNote:
		_, err := secretnote.Unmarshal(s.Data)
		if err != nil {
			return ErrInvalid
		}
		return nil
	default:
	}

	return ErrWrongType
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

	s.Data = encrypted
	s.IsEncrypted = true

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
	s.IsEncrypted = false

	return nil
}
