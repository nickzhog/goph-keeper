package cli

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/nickzhog/goph-keeper/internal/client/api"
	"github.com/nickzhog/goph-keeper/pkg/logging"
	"github.com/nickzhog/goph-keeper/pkg/secrets"
)

type cli struct {
	keeper api.KeeperClient
	logger *logging.Logger

	rlInstance       *readline.Instance
	passwordInputCfg *readline.Config

	isAuthenticated bool
	login           string

	secrets         map[int]secrets.AbstractSecret
	secretsViewMode bool
}

func New(keeper api.KeeperClient, logger *logging.Logger) *cli {
	c := &cli{
		keeper: keeper,
		logger: logger,
	}

	var err error
	c.rlInstance, err = readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		AutoComplete:    authCompleters,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}

	c.passwordInputCfg = c.rlInstance.GenPasswordConfig()
	c.passwordInputCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		c.rlInstance.SetPrompt(fmt.Sprintf("password(%v): ", len(line)))
		c.rlInstance.Refresh()
		return nil, 0, false
	})
	return c
}

func (c *cli) usage() {
	w := c.rlInstance.Stderr()

	io.WriteString(w, "commands:\n")
	if c.isAuthenticated {
		io.WriteString(w, secretsCompleters.Tree("    "))
	} else {
		io.WriteString(w, authCompleters.Tree("    "))
	}
}

func Start(ctx context.Context, keeper api.KeeperClient, logger *logging.Logger) {
	cli := New(keeper, logger)

	defer cli.rlInstance.Close()
	cli.rlInstance.CaptureExitSignal()

	log.SetOutput(cli.rlInstance.Stderr())

	cli.usage()

	for {
		line, err := cli.rlInstance.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch line {

		case helpCommand:
			cli.usage()

		case exitCommand:
			goto exit

		case "":

		default:
			if cli.isAuthenticated {
				cli.secretHandle(ctx, line)
			} else {
				cli.signHandle(ctx, line)
			}
		}
	}
exit:
}

func formatInput(s string) string {
	i := strings.Index(s, ":")
	s = s[i+1:]
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	return s
}
func formatNoteInput(s string) string {
	i := strings.Index(s, ":")
	s = s[i+1:]
	s = strings.TrimSpace(s)
	return s
}
