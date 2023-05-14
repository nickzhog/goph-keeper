package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func (c *client) ApplyToken(token string) {
	c.jwtTokenStr = token
}

func (c *client) AddTokenToMetadata(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("Authorization", "Bearer "+c.jwtTokenStr))
}

func (c *client) ResetToken() {
	c.jwtTokenStr = ""
}
