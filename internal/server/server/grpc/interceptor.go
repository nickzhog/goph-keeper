package grpc

import (
	"context"

	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(accountRep account.Repository, jwtSecretKey string) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if info.FullMethod == "/proto.Keeper/Register" ||
			info.FullMethod == "/proto.Keeper/Login" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokenMeta, ok := md["Authorization"]
		if !ok || len(tokenMeta) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		userID, err := account.ValidateJWT(tokenMeta[0], []byte(jwtSecretKey))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token invalid")
		}

		usr, err := accountRep.FindForID(ctx, userID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token invalid")
		}

		ctx = account.WriteUserToContext(ctx, usr)

		return handler(ctx, req)
	}
}
