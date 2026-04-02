package grpc

import (
	"context"
	pb "document-metadata/gen/proto/document"
	"document-metadata/pkg/document"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	mgr *document.Manager
	pb.UnimplementedDocumentServiceServer
}

func NewHandler(mgr *document.Manager) *Handler {
	return &Handler{mgr: mgr}
}

func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.Document, error) {
	doc, err := h.mgr.Create(ctx, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, document.ErrInvalidName) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to create document")
	}
	return &pb.Document{
		Id:          doc.ID,
		Name:        doc.Name,
		Description: doc.Description,
	}, nil
}

func (h *Handler) List(ctx context.Context, req *emptypb.Empty) (*pb.DocumentList, error) {
	doc, err := h.mgr.FindAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	pbDocs := []*pb.Document{}
	for _, d := range doc {
		pbdoc := pb.Document{
			CreatedAt:   d.CreatedAt.String(),
			Id:          d.ID,
			Name:        d.Name,
			Description: d.Description,
		}
		pbDocs = append(pbDocs, &pbdoc)
	}
	return &pb.DocumentList{
		Documents: pbDocs,
	}, nil
}

func (h *Handler) Get(ctx context.Context, req *pb.IDRequest) (*pb.Document, error) {
	doc, err := h.mgr.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "document not found")
	}
	return &pb.Document{
		CreatedAt:   doc.CreatedAt.String(),
		Id:          doc.ID,
		Name:        doc.Name,
		Description: doc.Description,
	}, nil
}

func (h *Handler) Delete(ctx context.Context, req *pb.IDRequest) (*emptypb.Empty, error) {
	if err := h.mgr.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.NotFound, "document not found")
	}
	return &emptypb.Empty{}, nil
}
