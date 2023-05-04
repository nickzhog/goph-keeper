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

	secret, err := k.srv.FindSecretByID(ctx, in.Secretid)
	if err != nil {
		return nil, err
	}

	if secret.UserID != usr.ID {
		return nil, status.Error(codes.Unauthenticated, "not your secret")
	}

	return &pb.GetSecretResponse{Secret: secret.Data}, nil
}
