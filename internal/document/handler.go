package document

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Process the request body and create a document.
	var document Document
	if err := json.Unmarshal(body, &document); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to create the document.
	result, err := h.service.CreateDocument(ctx, SrvCreateDocumentParam{
		Document: document,
	})
	if err != nil {
		http.Error(w, "Failed to create document", http.StatusInternalServerError)
		return
	}

	// Send the response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Document created successfully",
		"id":      result.ID,
	})
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	// Get the request context.
	ctx := r.Context()

	// Read the request URL parameters (if any).
	id := r.URL.Query().Get("id")

	// Call the service to list the documents.
	result, err := h.service.ListDocument(ctx, SrvListDocumentParam{
		ID: func() int64 {
			if id == "" {
				return 0
			}
			// In a real application, you would want to handle the error properly if the ID is not a valid integer.
			parsedID, _ := strconv.ParseInt(id, 10, 64)
			return parsedID
		}(),
	})
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
