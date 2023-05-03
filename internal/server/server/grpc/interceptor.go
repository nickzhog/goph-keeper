package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nickzhog/goph-keeper/internal/server/service/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(accountRep account.Repository, jwtKey string) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		tokenMeta, ok := md["token"]
		if !ok || len(tokenMeta) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenMeta[0], claims, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return jwtKey, nil
		})

		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, fmt.Errorf("invalid token")
		}

		expiresTime, exist := claims["exp"].(int64)
		if !exist {
			return nil, fmt.Errorf("expires time missed")
		}
		if expiresTime < time.Now().Unix() {
			return nil, fmt.Errorf("token not allowed")
		}

		userID, exist := claims[account.ClaimsKeyUserID].(string)
		if !exist {
			return nil, fmt.Errorf("user id missed")
		}

		usr, err := accountRep.FindForID(ctx, userID)
		if err != nil {
			return nil, err
		}

		ctx = account.WriteUserToContext(usr, ctx)

		return handler(ctx, req)
	}
}
