package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatInput(t *testing.T) {
	a := assert.New(t)
	s := "login:  input "
	new := formatInput(s)
	a.Equal("input", new)

	s = "input"
	new = formatInput(s)
	a.Equal("input", new)
}

func Test_formatNoteInput(t *testing.T) {
	a := assert.New(t)
	s := "note:  note note "
	new := formatNoteInput(s)
	a.Equal("note note", new)

	s = "text text 123 note"
	new = formatNoteInput(s)
	a.Equal("text text 123 note", new)
}
