package grpc

import (
	"context"

	pb "github.com/nickzhog/goph-keeper/internal/proto"
	"github.com/nickzhog/goph-keeper/internal/server/server"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"github.com/nickzhog/goph-keeper/internal/server/service/secrets"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.KeeperServer = (*KeeperServer)(nil)

type KeeperServer struct {
	srv server.Server

	pb.UnimplementedKeeperServer
}

func NewKeeperServer(srv server.Server) *KeeperServer {
	return &KeeperServer{
		srv: srv,
	}
}

func (k *KeeperServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := k.srv.Register(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RegisterResponse{Ok: true}, nil
}

func (k *KeeperServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse

	tokenStr, err := k.srv.Login(ctx, in.Login, in.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response.Token = tokenStr
	return &response, nil
}

func (k *KeeperServer) SecretsView(ctx context.Context, in *pb.SecretViewRequest) (*pb.SecretViewResponse, error) {
	usr := account.ReadUserFromContext(ctx)

	data, err := k.srv.FindSecretsForUser(ctx, usr.ID)
	if err == secrets.ErrNotFound {
		return &pb.SecretViewResponse{}, nil
	}

	var answer pb.SecretViewResponse

	for _, item := range data {
		answer.Secrets = append(answer.Secrets, &pb.SecretView{
			Id:    item.ID,
			Title: item.Title,
			Stype: pb.SecretType(pb.SecretType_value[item.SType]),
		})
	}

	return &answer, nil
}

func (k *KeeperServer) GetSecret(ctx context.Context, in *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	usr := account.ReadUserFromContext(ctx)

	secret, err := k.srv.FindSecretByID(ctx, in.Secretid, usr.ID)
	if err != nil {
		return nil, err
	}

	return &pb.GetSecretResponse{
		Secret: &pb.Secret{
			Title: secret.Title,
			Stype: pb.SecretType(pb.SecretType_value[secret.SType]),
			Id:    secret.ID,
			Data:  secret.Data,
		},
	}, nil
}

func (k *KeeperServer) CreateSecret(ctx context.Context, in *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	usr := account.ReadUserFromContext(ctx)

	secret := secrets.NewSecret(
		in.Secret.Id,
		usr.ID, in.Secret.Title,
		in.Secret.GetStype().String(),
		in.Secret.Data)

	err := k.srv.CreateSecret(ctx, *secret)
	if err != nil {
		return nil, err
	}

	return &pb.CreateSecretResponse{Ok: true}, nil
}

func (k *KeeperServer) DeleteSecret(ctx context.Context, in *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	usr := account.ReadUserFromContext(ctx)

	err := k.srv.DeleteSecret(ctx, in.Secretid, usr.ID)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteSecretResponse{Ok: true}, nil
}

func (k *KeeperServer) UpdateSecret(ctx context.Context, in *pb.UpdateSecretRequest) (*pb.UpdateSecretResponse, error) {
	usr := account.ReadUserFromContext(ctx)

	secret := secrets.NewSecret(
		in.Secret.Id,
		usr.ID, in.Secret.Title,
		in.Secret.GetStype().String(),
		in.Secret.Data)

	err := k.srv.UpdateSecret(ctx, *secret, usr.ID)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSecretResponse{Ok: true}, nil
}
