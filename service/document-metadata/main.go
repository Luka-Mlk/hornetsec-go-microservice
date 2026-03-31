package main

import (
	"context"
	"document-metadata/pkg/db/memory"
	"document-metadata/pkg/document"
	"document-metadata/pkg/transport/grpc"
	"document-metadata/pkg/transport/rest"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/k0kubun/pp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	memStore := memory.NewDocumentStore()
	mgr := document.NewManager(memStore)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// HTTP server
	wg.Add(1)
	go rest.Run(ctx, &wg, errChan, mgr)

	// gRPC server
	wg.Add(1)
	go grpc.Run(ctx, &wg, errChan, mgr)

	// Gracefull stop
	select {
	case err := <-errChan:
		pp.Printf("Server incorrectly started: %v\n", err)
		stop()
	case <-ctx.Done():
		pp.Println("Shutting down")
	}
	wg.Wait()
}
