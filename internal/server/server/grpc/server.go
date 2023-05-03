package grpc

import (
	"context"
	"net"

	pb "github.com/nickzhog/goph-keeper/internal/proto"
	"github.com/nickzhog/goph-keeper/internal/server/config"
	"github.com/nickzhog/goph-keeper/internal/server/server"
	"google.golang.org/grpc"
)

func Serve(ctx context.Context, srv server.Server, cfg *config.Config) {
	gRPCsrv := grpc.NewServer()

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
