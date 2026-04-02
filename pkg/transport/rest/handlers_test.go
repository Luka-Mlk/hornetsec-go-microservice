package rest_test

import (
	"bytes"
	"context"
	"document-metadata/pkg/db/memory"
	"document-metadata/pkg/document"
	"document-metadata/pkg/transport/rest"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() (*http.ServeMux, *document.Manager) {
	mgr := document.NewManager(memory.NewDocumentStore())
	h := rest.NewHandler(mgr)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux, mgr
}

func TestHandlers_AllEndpoints(t *testing.T) {
	mux, mgr := setup()
	ctx := context.Background()

	t.Run("POST /api/v1/documents - Success", func(t *testing.T) {
		body := `{"name": "Spec", "description": "Tech Doc"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", rr.Code)
		}
	})

	t.Run("POST /api/v1/documents - Validation Failure", func(t *testing.T) {
		body := `{"name": ""}` // Empty name triggers ErrInvalidName
		req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", bytes.NewBufferString(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("GET /api/v1/documents - List", func(t *testing.T) {
		mgr.Create(ctx, "Doc 1", "Desc")
		mgr.Create(ctx, "Doc 2", "Desc")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/documents", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		var docs []document.Document
		json.NewDecoder(rr.Body).Decode(&docs)

		if len(docs) < 2 {
			t.Errorf("expected at least 2 docs, got %d", len(docs))
		}
	})

	t.Run("GET /api/v1/documents/{id} - Success", func(t *testing.T) {
		doc, _ := mgr.Create(ctx, "GetMe", "Desc")

		path := fmt.Sprintf("/api/v1/documents/%s", doc.ID)
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})

	t.Run("DELETE /api/v1/documents/{id} - Success", func(t *testing.T) {
		doc, _ := mgr.Create(ctx, "DeleteMe", "Desc")

		path := fmt.Sprintf("/api/v1/documents/%s", doc.ID)
		req := httptest.NewRequest(http.MethodDelete, path, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rr.Code)
		}

		_, err := mgr.Get(ctx, doc.ID)
		if err == nil {
			t.Error("document should have been deleted")
		}
	})
}
