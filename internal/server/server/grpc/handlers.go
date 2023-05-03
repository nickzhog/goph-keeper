package grpc

import (
	"context"

	pb "github.com/nickzhog/goph-keeper/internal/proto"
	"github.com/nickzhog/goph-keeper/internal/server/server"
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

}

func (k *KeeperServer) GetSecret(ctx context.Context, in *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {

}
