package rest

import (
	"document-metadata/pkg/constants"
	"document-metadata/pkg/document"
	"encoding/json"
	"net/http"
)

const (
	apiRoot         = "/api/v1/documents"
	apiDocumentByID = "/api/v1/documents/{id}"
)

// Keeping dto's in handler file for the sake of keeping things simple by go standard
type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MiddlewareFunc func(next http.Handler) http.Handler

type Handler struct {
	mgr *document.Manager
}

func NewHandler(mgr *document.Manager) *Handler {
	return &Handler{mgr: mgr}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST "+apiRoot, h.create)
	mux.HandleFunc("GET "+apiRoot, h.list)
	mux.HandleFunc("GET "+apiDocumentByID, h.get)
	mux.HandleFunc("DELETE "+apiDocumentByID, h.delete)
}

func (h *Handler) Use(mux *http.ServeMux, mf MiddlewareFunc) http.Handler {
	return mf(mux)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	doc, err := h.mgr.Create(req.Name, req.Description)
	if err != nil {
		http.Error(w, "failed to create document", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", constants.HeaderContentTypeApplicationJson)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	doc, err := h.mgr.FindAll()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", constants.HeaderContentTypeApplicationJson)
	json.NewEncoder(w).Encode(doc)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	doc, err := h.mgr.Get(id)
	if err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", constants.HeaderContentTypeApplicationJson)
	json.NewEncoder(w).Encode(doc)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.mgr.Delete(id); err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
