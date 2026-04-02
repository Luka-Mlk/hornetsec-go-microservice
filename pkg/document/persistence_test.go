package document_test

import (
	"context"
	"document-metadata/pkg/db/memory"
	"document-metadata/pkg/document"
	"testing"
)

func TestManager_Create(t *testing.T) {
	repo := memory.NewDocumentStore()
	mgr := document.NewManager(repo)
	ctx := context.Background()

	tests := []struct {
		name        string
		docName     string
		description string
		wantErr     bool
	}{
		{
			name:        "Success: Valid inputs",
			docName:     "Project Alpha",
			description: "Internal specifications",
			wantErr:     false,
		},
		{
			name:        "Failure: Empty name",
			docName:     "",
			description: "Should fail validation",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mgr.Create(ctx, tt.docName, tt.description)

			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Name != tt.docName {
					t.Errorf("expected name %s, got %s", tt.docName, got.Name)
				}
				if got.ID == "" {
					t.Error("expected a generated ID, got empty string")
				}
			}
		})
	}
}

func TestManager_List(t *testing.T) {
	repo := memory.NewDocumentStore()
	mgr := document.NewManager(repo)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		mgr.Create(ctx, "Doc", "Desc")
	}

	t.Run("Success: List all documents", func(t *testing.T) {
		docs, err := mgr.FindAll(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(docs) != 3 {
			t.Errorf("expected 3 documents, got %d", len(docs))
		}
	})
}

func TestManager_Get(t *testing.T) {
	repo := memory.NewDocumentStore()
	mgr := document.NewManager(repo)
	ctx := context.Background()

	seed, _ := mgr.Create(ctx, "Seed Doc", "Desc")

	t.Run("Success: Found existing document", func(t *testing.T) {
		got, err := mgr.Get(ctx, seed.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != seed.ID {
			t.Errorf("expected ID %s, got %s", seed.ID, got.ID)
		}
	})

	t.Run("Failure: Document not found", func(t *testing.T) {
		_, err := mgr.Get(ctx, "non-existent-id")
		if err == nil {
			t.Error("expected an error for missing ID, got nil")
		}
	})
}

func TestManager_Delete(t *testing.T) {
	repo := memory.NewDocumentStore()
	mgr := document.NewManager(repo)
	ctx := context.Background()

	seed, _ := mgr.Create(ctx, "To Delete", "Desc")

	t.Run("Success: Delete document", func(t *testing.T) {
		if err := mgr.Delete(ctx, seed.ID); err != nil {
			t.Fatalf("failed to delete: %v", err)
		}
		_, err := mgr.Get(ctx, seed.ID)
		if err == nil {
			t.Error("expected error getting deleted document, got nil")
		}
	})
}

func TestManager_ConcurrentCreate(t *testing.T) {
	repo := memory.NewDocumentStore()
	mgr := document.NewManager(repo)
	ctx := context.Background()

	const count = 100
	done := make(chan bool)

	for i := 0; i < count; i++ {
		go func(i int) {
			_, err := mgr.Create(ctx, "Doc", "Desc")
			if err != nil {
				t.Errorf("concurrent create failed: %v", err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < count; i++ {
		<-done
	}
}
