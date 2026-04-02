package grpc

import (
	"context"
	"document-metadata/pkg/constants"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var id string

	if ok && len(md[constants.HeaderRequestID]) > 0 {
		id = md[constants.HeaderRequestID][0]
	} else {
		id = uuid.NewString()
	}

	header := metadata.Pairs(constants.HeaderRequestID, id)
	grpc.SendHeader(ctx, header)

	newCtx := context.WithValue(ctx, constants.HeaderRequestID, id)
	return handler(newCtx, req)
}
