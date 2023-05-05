package secretnote

import (
	"encoding/json"
	"errors"
)

type SecretNote struct {
	Note string
}

func (n *SecretNote) IsValid() bool {
	if n.Note == "" {
		return false
	}

	return true
}

func Unmarshal(secretData []byte) (*SecretNote, error) {
	var note SecretNote
	err := json.Unmarshal(secretData, &note)
	if err != nil {
		return nil, err
	}

	if !note.IsValid() {
		return nil, errors.New("not valid secret")
	}

	return &note, nil
}
