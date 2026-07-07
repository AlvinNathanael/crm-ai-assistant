package document

import (
	"encoding/json"
	"errors"
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
	var param HandlerCreateDocumentParam
	if err := json.Unmarshal(body, &param); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(HandlerCreateDocumentResponse{
		Message: "Document created successfully",
		ID:      result.ID,
	})
	w.WriteHeader(http.StatusCreated)
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

	// Read the request URL (if any).
	id := r.PathValue("id")

	// Call the service to get the document.
	result, err := h.service.GetDocument(ctx, SrvGetDocumentParam{
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
