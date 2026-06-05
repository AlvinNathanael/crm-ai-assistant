package document

import "net/http"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request)

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request)
