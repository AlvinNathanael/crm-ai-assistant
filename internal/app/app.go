package app

import (
	"net/http"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
)

type App struct {
	Repos    *Repositories
	Services *Services
	Mux      *http.ServeMux
}

type Repositories struct {
	Document document.RepositoryItf
}

type Services struct {
	Document document.ServiceItf
}

func NewApp() *App {
	// Initialize repositories.
	repos := InitRepositories()
	// Initialize services.
	services := InitServices(repos)
	// Initialize the HTTP serve mux.
	mux := InitHandlers(services)

	return &App{
		Repos:    repos,
		Services: services,
		Mux:      mux,
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
