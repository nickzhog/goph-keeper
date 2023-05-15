package cli

import (
	"context"
	"log"
	"strconv"

	"github.com/nickzhog/goph-keeper/internal/client/service/auth"
)

func (c *cli) signHandle(ctx context.Context, line string) {
	switch line {
	case loginCommand:
		login, err := c.rlInstance.ReadlineWithDefault("your login: ")
		if err != nil {
			log.Println("error:", strconv.Quote(err.Error()))
			return
		}
		login = formatInput(login)

		pswd, err := c.rlInstance.ReadPasswordWithConfig(c.passwordInputCfg)
		if err != nil {
			log.Println("error:", strconv.Quote(err.Error()))
			return
		}
		err = auth.Login(ctx, c.keeper, login, string(pswd))
		if err != nil {
			log.Println("wrong password")
			return
		}
		c.isAuthenticated = true
		c.login = login
		log.Println("successful,", c.login)
		c.usage()

	case registerCommand:
		login, err := c.rlInstance.ReadlineWithDefault("your login: ")
		if err != nil {
			log.Println("error:", strconv.Quote(err.Error()))
			return
		}
		login = formatInput(login)

		pswd, err := c.rlInstance.ReadPasswordWithConfig(c.passwordInputCfg)
		if err != nil {
			log.Println("error:", strconv.Quote(err.Error()))
			return
		}

		pswdRepeat, err := c.rlInstance.ReadPasswordWithConfig(c.passwordInputCfg)
		if err != nil {
			log.Println("error:", strconv.Quote(err.Error()))
			return
		}
		if string(pswd) != string(pswdRepeat) {
			log.Println("passwords not equal")
			return
		}

		err = auth.Register(ctx, c.keeper, login, string(pswd))
		if err != nil {
			log.Println("invalid login or password:", err.Error())
			return
		}
		log.Println("registration complete, log in please")

	default:
		log.Println("unkown command:", strconv.Quote(line))
		c.usage()

	}
}
