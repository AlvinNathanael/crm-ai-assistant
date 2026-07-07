package app

import (
	"database/sql"
	"net/http"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/config"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/database"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
)

type App struct {
	DB  *sql.DB
	Mux *http.ServeMux

	Repos    *Repositories
	Services *Services
}

type Repositories struct {
	Document document.RepositoryItf
}

type Services struct {
	Document document.ServiceItf
}

func NewApp(cfg config.Config) (*App, error) {
	// Init DB.
	db, err := database.New(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	// Initialize repositories.
	repos := InitRepositories(db)

	// Initialize services.
	services := InitServices(repos)

	// Initialize the HTTP serve mux.
	mux := InitHandlers(services)

	return &App{
		DB:       db,
		Repos:    repos,
		Services: services,
		Mux:      mux,
	}, nil
}

func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Document: document.NewRepository(db),
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
	// TODO: move route registration to feature handlers as routes grow.
	mux.HandleFunc("POST /documents", handlerDocument.CreateDocument)
	mux.HandleFunc("GET /documents", handlerDocument.ListDocuments)
	mux.HandleFunc("GET /documents/{id}", handlerDocument.GetDocument)

	return mux
}
