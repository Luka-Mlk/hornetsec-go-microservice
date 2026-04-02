package grpc

import (
	"context"
	"document-metadata/pkg/document"
	"net"
	"os"
	"sync"

	pb "document-metadata/gen/proto/document"

	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	protocol = "tcp"
)

func Run(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error, mgr *document.Manager) {
	defer wg.Done()
	addr := os.Getenv("GRPC_PORT")
	socket, err := net.Listen(protocol, addr)
	if err != nil {
		errChan <- err
		return
	}
	handler := NewHandler(mgr)
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(RequestIDInterceptor),
	)
	pb.RegisterDocumentServiceServer(srv, handler)
	reflection.Register(srv)
	pp.Println("[GRPC] started")
	go func() {
		<-ctx.Done()
		pp.Println("[GRPC] stopping")
		srv.GracefulStop()
	}()
	if err := srv.Serve(socket); err != nil {
		errChan <- err
	}
}
