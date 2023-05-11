package secretnote

import (
	"encoding/json"
	"errors"
)

type SecretNote struct {
	Note string `json:"note,omitempty"`
}

func (n *SecretNote) IsValid() bool {
	if n.Note == "" {
		return false
	}

	return true
}

func (n *SecretNote) Marshal() []byte {
	data, _ := json.Marshal(n)
	return data
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

func New(note string) *SecretNote {
	return &SecretNote{
		Note: note,
	}
}
