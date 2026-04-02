package rest

import (
	"context"
	"document-metadata/pkg/document"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/k0kubun/pp"
)

func Run(ctx context.Context, wg *sync.WaitGroup, errChan chan<- error, mgr *document.Manager) {
	defer wg.Done()
	mux := http.NewServeMux()
	h := NewHandler(mgr)
	h.RegisterRoutes(mux)

	handlers := h.Use(mux, WithRequestID)
	addr := os.Getenv("HTTP_PORT")

	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: handlers,
	}

	pp.Println("[REST] started")
	go func() {
		<-ctx.Done()
		log.Println("[REST] stopping")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
	}
}
