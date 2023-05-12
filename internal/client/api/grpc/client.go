package grpc

import (
	"context"
	"log"

	pb "github.com/nickzhog/goph-keeper/api/proto"
	"github.com/nickzhog/goph-keeper/internal/client/api"
	"github.com/nickzhog/goph-keeper/pkg/logging"
	"github.com/nickzhog/goph-keeper/pkg/secrets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var _ api.KeeperClient = (*client)(nil)

type client struct {
	c           pb.KeeperClient
	jwtTokenStr string
	logger      *logging.Logger
}

func (c *client) ApplyJWT(token string) {
	c.jwtTokenStr = token
}

func (c *client) AddTokenToMetadata(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("Authorization", "Bearer "+c.jwtTokenStr))
}

func NewClient(addr, certPath string, logger *logging.Logger) *client {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}

	c := pb.NewKeeperClient(conn)

	return &client{logger: logger, c: c}
}

func (c *client) CreateAccount(ctx context.Context, login string, password string) error {
	_, err := c.c.Register(ctx, &pb.RegisterRequest{Login: login, Password: password})
	return err
}

func (c *client) LoginAccount(ctx context.Context, login string, password string) (token string, err error) {
	response, err := c.c.Login(ctx, &pb.LoginRequest{Login: login, Password: password})
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *client) SecretsView(ctx context.Context) ([]secrets.AbstractSecret, error) {

	ctx = c.AddTokenToMetadata(ctx)

	response, err := c.c.SecretsView(ctx, &pb.SecretViewRequest{})

	if err != nil {
		return nil, err
	}

	secretsView := make([]secrets.AbstractSecret, 0, len(response.Secrets))

	for _, s := range response.Secrets {
		secretsView = append(secretsView, *secrets.NewSecretWithoutEncryptedData(
			s.Id,
			"",
			s.Title,
			s.Stype.String(),
			nil,
		))
	}
	return secretsView, nil
}

func (c *client) CreateSecret(ctx context.Context, secret secrets.AbstractSecret) error {
	ctx = c.AddTokenToMetadata(ctx)

	_, err := c.c.CreateSecret(ctx, &pb.CreateSecretRequest{
		Secret: &pb.Secret{
			Id:    secret.ID,
			Title: secret.Title,
			Stype: pb.SecretType(pb.SecretType_value[secret.SType]),
			Data:  secret.Data,
		},
	})

	return err
}

func (c *client) GetSecretByID(ctx context.Context, secretID string) (secrets.AbstractSecret, error) {
	ctx = c.AddTokenToMetadata(ctx)

	response, err := c.c.GetSecret(ctx, &pb.GetSecretRequest{
		Secretid: secretID,
	})
	if err != nil {
		return secrets.AbstractSecret{}, err
	}

	return *secrets.NewSecretWithoutEncryptedData(
		response.Secret.Id,
		"",
		response.Secret.Title,
		response.Secret.Stype.String(),
		response.Secret.Data,
	), nil
}

func (c *client) UpdateSecretByID(ctx context.Context, secret secrets.AbstractSecret) error {
	ctx = c.AddTokenToMetadata(ctx)

	_, err := c.c.UpdateSecret(ctx, &pb.UpdateSecretRequest{
		Secret: &pb.Secret{
			Id:    secret.ID,
			Title: secret.Title,
			Stype: pb.SecretType(pb.SecretType_value[secret.SType]),
			Data:  secret.Data,
		},
	})

	return err
}

func (c *client) DeleteSecretByID(ctx context.Context, secretID string) error {
	ctx = c.AddTokenToMetadata(ctx)

	_, err := c.c.DeleteSecret(ctx, &pb.DeleteSecretRequest{Secretid: secretID})
	return err
}
