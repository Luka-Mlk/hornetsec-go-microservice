package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var id string

	if ok && len(md["X-Request-ID"]) > 0 {
		id = md["X-Request-ID"][0]
	} else {
		id = uuid.NewString()
	}

	header := metadata.Pairs("X-Request-ID", id)
	grpc.SendHeader(ctx, header)

	newCtx := context.WithValue(ctx, "X-Request-ID", id)
	return handler(newCtx, req)
}
