package secretbinary

import (
	"encoding/json"
	"errors"
)

type SecretBinary struct {
	Data []byte
	Note string
}

func (b *SecretBinary) IsValid() bool {
	if len(b.Data) == 0 && b.Note == "" {
		return false
	}

	return true
}

func Unmarshal(secretData []byte) (*SecretBinary, error) {
	var data SecretBinary
	err := json.Unmarshal(secretData, &data)
	if err != nil {
		return nil, err
	}

	if !data.IsValid() {
		return nil, errors.New("not valid secret")
	}

	return &data, nil
}
