package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nickzhog/goph-keeper/internal/client/api/grpc"
	"github.com/nickzhog/goph-keeper/internal/client/config"
	"github.com/nickzhog/goph-keeper/pkg/logging"
	"github.com/nickzhog/goph-keeper/pkg/secrets"
	secretaccount "github.com/nickzhog/goph-keeper/pkg/secrets/account"
)

func main() {
	cfg := config.GetConfig()

	ctx := context.Background()
	logger := logging.GetOnlyFileLogger()
	keeperClient := grpc.NewClient(cfg.Settings.ConnAddress, cfg.Settings.CertSSL, logger)

	fmt.Println("sign up:")
	fmt.Println(keeperClient.CreateAccount(ctx, "user", "password"))

	////////////////////////////////////////////////////////////
	fmt.Println("sign in:")
	jwtTokenStr, err := keeperClient.LoginAccount(ctx, "user", "password")
	keeperClient.ApplyJWT(jwtTokenStr)
	fmt.Println(jwtTokenStr, err)

	////////////////////////////////////////////////////////////
	fmt.Println("create secret:")
	sAccount := secretaccount.New("site.com", "login", "password", "", "note")
	data := sAccount.Marshal()
	sAbstract := secrets.NewSecretWithoutEncryptedData("secret_id", "user_id", "title", secrets.TypeAccount, data)

	fmt.Println(keeperClient.CreateSecret(ctx, *sAbstract))

	////////////////////////////////////////////////////////////
	fmt.Println("secret view:")
	secretsview, err := keeperClient.SecretsView(ctx)

	fmt.Println("err:", err)
	for i, item := range secretsview {
		fmt.Println(i+1, item.Title, item.SType)
	}
	if len(secretsview) < 1 {
		log.Fatal("secretsview < 1")
	}

	////////////////////////////////////////////////////////////
	fmt.Println("get one secret:")
	s, err := keeperClient.GetSecretByID(ctx, secretsview[0].ID)
	fmt.Println("err:", err)
	fmt.Printf("secret: %+v\n", s)

	////////////////////////////////////////////////////////////
	fmt.Println("update secret:")
	s.Title += "_test"
	fmt.Println("err:", keeperClient.UpdateSecretByID(ctx, s))

	s, err = keeperClient.GetSecretByID(ctx, secretsview[0].ID)
	fmt.Println("find err:", err)
	fmt.Println("new title:", s.Title)

	////////////////////////////////////////////////////////////
	fmt.Println("delete secret:")
	fmt.Println("err:", keeperClient.DeleteSecretByID(ctx, s.ID))
	s, err = keeperClient.GetSecretByID(ctx, s.ID)
	if err == nil {
		fmt.Println("ERROR: can find deleted secret")
	}
}
