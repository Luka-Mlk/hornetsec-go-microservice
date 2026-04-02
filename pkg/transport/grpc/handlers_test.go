package grpc_test

import (
	"context"
	pb "document-metadata/gen/proto/document"
	"document-metadata/pkg/db/memory"
	"document-metadata/pkg/document"
	transport "document-metadata/pkg/transport/grpc"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

func setup(t *testing.T) (pb.DocumentServiceClient, *document.Manager) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	mgr := document.NewManager(memory.NewDocumentStore())
	s := grpc.NewServer()
	pb.RegisterDocumentServiceServer(s, transport.NewHandler(mgr))

	go func() {
		if err := s.Serve(lis); err != nil {
			return
		}
	}()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}

	t.Cleanup(func() {
		conn.Close()
		s.Stop()
	})

	return pb.NewDocumentServiceClient(conn), mgr
}

func TestGRPCHandlers(t *testing.T) {
	client, mgr := setup(t)
	ctx := context.Background()

	t.Run("Create - Success", func(t *testing.T) {
		resp, err := client.Create(ctx, &pb.CreateRequest{Name: "GRPC Doc"})
		if err != nil || resp.Name != "GRPC Doc" {
			t.Errorf("failed to create via gRPC: %v", err)
		}
	})

	t.Run("Create - Invalid Argument", func(t *testing.T) {
		_, err := client.Create(ctx, &pb.CreateRequest{Name: ""})
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.InvalidArgument {
			t.Errorf("expected InvalidArgument, got %v", st.Code())
		}
	})

	t.Run("Get - Not Found", func(t *testing.T) {
		_, err := client.Get(ctx, &pb.IDRequest{Id: "missing"})
		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Errorf("expected NotFound, got %v", st.Code())
		}
	})

	t.Run("List - Success", func(t *testing.T) {
		mgr.Create(ctx, "Doc 1", "")
		resp, err := client.List(ctx, &emptypb.Empty{})
		if err != nil || len(resp.Documents) == 0 {
			t.Errorf("list failed: %v", err)
		}
	})
}
