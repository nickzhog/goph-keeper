package cli

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/nickzhog/goph-keeper/internal/client/api"
)

func (c *cli) secretHandle(ctx context.Context, line string) {
	switch {
	case line == secretsCommand:
		secretsView, err := c.keeper.SecretsView(ctx)
		if err != nil {
			if errors.Is(err, api.ErrInvalidToken) {
				c.logout()
				log.Println("need to auth")
				break
			}

			log.Println("error: ", err.Error())
			break
		}

		if len(secretsView) < 1 {
			log.Println("you dont have secrets")
			break
		}
		c.showSecrets(secretsView)

	case line == createCommand:
		c.createSecret(ctx)

	case strings.HasPrefix(line, getCommand):
		if !c.secretsViewMode {
			log.Println("request secrets first")
			break
		}
		if len(line) < len(getCommand)+2 {
			log.Println(getCommand, "<number>")
			break
		}

		numS := line[len(getCommand)+1:]
		num, err := strconv.Atoi(numS)
		if err != nil {
			log.Println("wrong number")
			break
		}
		secret, exist := c.secrets[num]
		if !exist {
			log.Println("wrong number")
			break
		}
		c.getSecret(ctx, secret.ID)

	// case strings.HasPrefix(line, updateCommand):

	case strings.HasPrefix(line, deleteCommand):
		if !c.secretsViewMode {
			log.Println("request secrets first")
			break
		}
		if len(line) < len(deleteCommand)+2 {
			log.Println(deleteCommand, "<number>")
			break
		}
		numS := line[len(deleteCommand)+1:]
		num, err := strconv.Atoi(numS)
		if err != nil {
			log.Println("wrong number")
			break
		}
		secret, exist := c.secrets[num]
		if !exist {
			log.Println("wrong number")
			break
		}
		err = c.keeper.DeleteSecretByID(ctx, secret.ID)
		if err != nil {
			log.Println("error:", err.Error())
		}
		log.Println("deleted")

	default:
		log.Println("unkown command:", strconv.Quote(line))
		c.usage()

	}
}

func (c *cli) logout() {
	c.secrets = nil
	c.secretsViewMode = false
	c.keeper.ResetToken()
	c.isAuthenticated = false
	c.login = ""
}
