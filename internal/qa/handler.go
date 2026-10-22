package qa

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Ask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var request struct {
		Question string `json:"question"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(request.Question) == "" {
		http.Error(w, "Question is required", http.StatusBadRequest)
		return
	}

	result, err := h.service.Ask(r.Context(), request.Question)
	if err != nil {
		if errors.Is(err, ErrUpstream) {
			http.Error(w, "AI provider request failed", http.StatusBadGateway)
		} else {
			http.Error(w, "Failed to search document chunks", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
