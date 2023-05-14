package cli

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nickzhog/goph-keeper/pkg/secrets"
	secretaccount "github.com/nickzhog/goph-keeper/pkg/secrets/account"
	secretbinary "github.com/nickzhog/goph-keeper/pkg/secrets/binary"
	secretcard "github.com/nickzhog/goph-keeper/pkg/secrets/card"
	secretnote "github.com/nickzhog/goph-keeper/pkg/secrets/note"
)

func (c *cli) showSecrets(view []secrets.AbstractSecret) {
	c.secretsViewMode = true
	c.secrets = make(map[int]secrets.AbstractSecret, len(view))
	for i, item := range view {
		c.secrets[i+1] = item
		fmt.Printf("%d : %s : %s\n",
			len(c.secrets), item.Title, strings.ToLower(item.SType))
	}
	fmt.Printf("you have %d secrets\n", len(view))
	fmt.Printf("type '%s <number>' to get secret data\n", getCommand)
	// log.Printf("type '%s <number>' to update secret data\n", updateCommand)
	fmt.Printf("type '%s <number>' to delete secret\n", deleteCommand)
}

func (c *cli) getSecret(ctx context.Context, secretID string) {
	s, err := c.keeper.GetSecretByID(ctx, secretID)
	if err != nil {
		log.Println("error:", err.Error())
		return
	}

	switch s.SType {
	case secrets.TypeAccount:
		ac, err := secretaccount.Unmarshal(s.Data)
		if err != nil {
			log.Println("error with get secret", err.Error())
			return
		}
		ac.TotpCheck()
		fmt.Printf("domain: %s\nlogin: %s\npassword: %s\ntotp code: %s\nnote: %s",
			ac.SiteDomain, ac.Login, ac.Password, ac.CodeTOTP, ac.Note)

	case secrets.TypeBinary:
		fileName, err := c.rlInstance.ReadlineWithDefault("file name for save: ")
		if err != nil {
			fmt.Println("error:", strconv.Quote(err.Error()))
			return
		}
		fileName = formatInput(fileName)
		bin, err := secretbinary.Unmarshal(s.Data)
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		err = bin.SaveToFile(fileName)
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		fmt.Println("file saved")

	case secrets.TypeCard:
		card, err := secretcard.Unmarshal(s.Data)
		if err != nil {
			log.Println("error with get secret", err.Error())
			return
		}
		fmt.Printf("number: %s\nmonth: %s\nyear: %s\ncvv: %s\nholder name: %s\nnote: %s",
			card.Number, card.Month, card.Year, card.CVV, card.HolderName, card.Note)
	case secrets.TypeNote:
		note, err := secretnote.Unmarshal(s.Data)
		if err != nil {
			log.Println("error with get secret", err.Error())
			return
		}
		fmt.Printf("note: %s\n", note.Note)

	default:
		log.Println("error with get secret")
	}
}
