package main

import (
	"log"
	"net/http"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
)

type Repositories struct {
	Document *document.Repository
}

type Services struct {
	Document *document.Service
}

func main() {
	// Initialize repositories.
	repos := InitRepositories()

	// Initialize services.
	services := InitServices(repos)

	// Initialize handlers.
	mux := InitHandlers(services)

	// Listen and serve.
	log.Println("Server starting on: 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln("Failed HTTP server: " + err.Error())
	}
}

func InitRepositories() *Repositories {
	return &Repositories{
		Document: document.NewRepository(),
	}
}

func InitServices(repos *Repositories) *Services {
	return &Services{
		Document: document.NewService(repos.Document),
	}
}

func InitHandlers(services *Services) *http.ServeMux {
	// Initalize mux server.
	mux := http.NewServeMux()

	// Health check.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Initialize handlers.
	handlerDocument := document.NewHandler(services.Document)

	// Initialize handler paths.
	mux.HandleFunc("POST /documents", handlerDocument.CreateDocument)
	mux.HandleFunc("GET /documents", handlerDocument.ListDocuments)

	return mux
}
