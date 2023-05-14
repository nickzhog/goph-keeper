package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nickzhog/goph-keeper/pkg/secrets"
	secretaccount "github.com/nickzhog/goph-keeper/pkg/secrets/account"
	secretbinary "github.com/nickzhog/goph-keeper/pkg/secrets/binary"
	secretcard "github.com/nickzhog/goph-keeper/pkg/secrets/card"
	secretnote "github.com/nickzhog/goph-keeper/pkg/secrets/note"
)

func (c *cli) createSecret(ctx context.Context) {
	title, err := c.rlInstance.ReadlineWithDefault("secret title: ")
	if err != nil {
		fmt.Println("error:", strconv.Quote(err.Error()))
		return
	}
	title = formatInput(title)

	fmt.Printf("1 - account\n2 - binary\n3 - note\n4 - card\n")
	stypeNum, err := c.rlInstance.ReadlineWithDefault("secret type (1-4): ")
	if err != nil {
		fmt.Println("error:", strconv.Quote(err.Error()))
		return
	}
	stypeNum = strings.TrimSpace(stypeNum)
	stypeNum = stypeNum[len(stypeNum)-2:]
	stypeNum = strings.TrimSpace(stypeNum)

	var secretdata []byte
	var stype string
	switch stypeNum {
	case "1":
		stype = secrets.TypeAccount
		sitedomain, err := c.rlInstance.ReadlineWithDefault("site domain: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		sitedomain = formatInput(sitedomain)

		login, err := c.rlInstance.ReadlineWithDefault("login: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		login = formatInput(login)

		password, err := c.rlInstance.ReadPasswordWithConfig(c.passwordInputCfg)
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}

		keytotp, err := c.rlInstance.ReadlineWithDefault("totp key: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		keytotp = formatInput(keytotp)

		note, err := c.rlInstance.ReadlineWithDefault("note: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		note = formatInput(note)

		account := secretaccount.New(sitedomain, login, string(password), keytotp, note)
		if !account.IsValid() {
			fmt.Println("secret invalid")
			return
		}
		secretdata = account.Marshal()
	case "2":
		stype = secrets.TypeBinary
		filePath, err := c.rlInstance.ReadlineWithDefault("file path: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		filePath = formatInput(filePath)

		note, err := c.rlInstance.ReadlineWithDefault("note: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		note = formatInput(note)

		bin, err := secretbinary.New(filePath, note)
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		secretdata = bin.Marshal()

	case "3":
		stype = secrets.TypeNote
		note, err := c.rlInstance.ReadlineWithDefault("note: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		note = formatInput(note)

		snote := secretnote.New(note)
		secretdata = snote.Marshal()

	case "4":
		stype = secrets.TypeCard
		number, err := c.rlInstance.ReadlineWithDefault("number: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		number = formatInput(number)

		month, err := c.rlInstance.ReadlineWithDefault("month: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		month = formatInput(month)

		year, err := c.rlInstance.ReadlineWithDefault("year: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		year = formatInput(year)

		cvv, err := c.rlInstance.ReadlineWithDefault("cvv: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		cvv = formatInput(cvv)
		holder, err := c.rlInstance.ReadlineWithDefault("holder name: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		holder = formatInput(holder)

		note, err := c.rlInstance.ReadlineWithDefault("note: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		note = formatInput(note)

		card := secretcard.New(number, month, year, cvv, holder, note)
		if !card.IsValid() {
			fmt.Println("secret invalid")
			return
		}

		secretdata = card.Marshal()
	default:
		fmt.Println("wrong type", stypeNum)
		return
	}

	s := secrets.NewSecretWithoutEncryptedData("", "", title, stype, secretdata)
	err = c.keeper.CreateSecret(ctx, *s)
	if err != nil {
		fmt.Println("error:", strconv.Quote(err.Error()))
		return
	}
	fmt.Println("secret created")
}
