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
