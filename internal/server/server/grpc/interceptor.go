package grpc

import (
	"context"
	"strings"

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

		tokenMeta := md.Get("Authorization")
		if len(tokenMeta) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		tokenParts := strings.Split(tokenMeta[0], " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		token := tokenParts[1]

		userID, err := account.ValidateJWT(token, []byte(jwtSecretKey))
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
