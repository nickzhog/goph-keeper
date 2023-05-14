package secretbinary

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

func (b *SecretBinary) SaveToFile(fileName string) error {
	return os.WriteFile(fileName, b.Data, 0755)
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

func New(filepath, note string) (*SecretBinary, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	// 1 MB is max
	if info.Size() > 1024*1024 {
		return nil, fmt.Errorf("file size exceeds 1MB")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return &SecretBinary{
		Data: data,
		Note: note,
	}, nil
}
