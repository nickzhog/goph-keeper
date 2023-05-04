package grpc

import (
	"context"
	"net"

	pb "github.com/nickzhog/goph-keeper/internal/proto"
	"github.com/nickzhog/goph-keeper/internal/server/config"
	"github.com/nickzhog/goph-keeper/internal/server/server"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Serve(ctx context.Context, srv server.Server, cfg *config.Config, accountRep account.Repository) {
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		srv.Logger.Fatal(err)
	}

	gRPCsrv := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(NewAuthInterceptor(accountRep, cfg.Settings.JWTkey)),
	)

	pb.RegisterKeeperServer(gRPCsrv, NewKeeperServer(srv))
	go func() {
		listen, err := net.Listen("tcp", cfg.Settings.AddressGRPC)
		if err != nil {
			srv.Logger.Fatal(err)
		}
		if err = gRPCsrv.Serve(listen); err != nil && err != grpc.ErrServerStopped {
			srv.Logger.Fatalf("grpc listen:%+s\n", err)
		}
	}()

	srv.Logger.Tracef("grpc server started")

	<-ctx.Done()

	gRPCsrv.GracefulStop()

	srv.Logger.Tracef("grpc server stopped")
}
