package document

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service ServiceItf
}

func NewHandler(service ServiceItf) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Get the request context.
	ctx := r.Context()

	// Read the request body.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 10<<20))
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Process the request body and create a document.
	var param HandlerCreateDocumentParam
	if err := json.Unmarshal(body, &param); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(param.Title) == "" || strings.TrimSpace(param.Content) == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Call the service to create the document.
	result, err := h.service.CreateDocument(ctx, SrvCreateDocumentParam{
		Title:   param.Title,
		Content: param.Content,
	})
	if err != nil {
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	// Send the response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(HandlerCreateDocumentResponse{
		Message: "Document created successfully",
		ID:      result.ID,
	})
}

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	// Get the request context.
	ctx := r.Context()

	// Call the service to list the documents.
	result, err := h.service.ListDocuments(ctx, SrvListDocumentsParam{})
	if err != nil {
		http.Error(w, "Failed to list documents", http.StatusInternalServerError)
		return
	}

	// Send the response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"documents": result.Documents,
	})
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	// Get the request context.
	ctx := r.Context()

	id, ok := parseDocumentID(w, r)
	if !ok {
		return
	}

	// Call the service to get the document.
	result, err := h.service.GetDocument(ctx, SrvGetDocumentParam{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get document", http.StatusInternalServerError)
		}
		return
	}

	// Send the response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"document": result.Document,
	})
}

func (h *Handler) ReindexDocument(w http.ResponseWriter, r *http.Request) {
	id, ok := parseDocumentID(w, r)
	if !ok {
		return
	}

	result, err := h.service.ReindexDocument(r.Context(), SrvReindexDocumentParam{ID: id})
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to reindex document", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func parseDocumentID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Document ID must be a positive integer", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
