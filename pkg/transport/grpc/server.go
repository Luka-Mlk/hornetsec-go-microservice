package grpc

import (
	"context"
	"document-metadata/pkg/document"
	"net"
	"sync"

	pb "document-metadata/gen/proto/document"

	"github.com/k0kubun/pp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error, mgr *document.Manager) {
	defer wg.Done()
	socket, err := net.Listen("tcp", ":9090")
	if err != nil {
		errChan <- err
		return
	}
	handler := NewHandler(mgr)
	srv := grpc.NewServer()
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
