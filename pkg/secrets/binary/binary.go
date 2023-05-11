package secretbinary

import (
	"encoding/json"
	"errors"
)

type SecretBinary struct {
	Data []byte `json:"data,omitempty"`
	Note string `json:"note,omitempty"`
}

func (b *SecretBinary) IsValid() bool {
	if len(b.Data) == 0 && b.Note == "" {
		return false
	}

	return true
}

func (b *SecretBinary) Marshal() []byte {
	data, _ := json.Marshal(b)
	return data
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

func New(data []byte, note string) *SecretBinary {
	return &SecretBinary{
		Data: data,
		Note: note,
	}
}
