package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/config"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/database"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/embedding"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/qa"
)

type App struct {
	DB  *sql.DB
	Mux *http.ServeMux

	Repos    *Repositories
	Services *Services
	QA       *qa.Service
}

type Repositories struct {
	Document document.RepositoryItf
}

type Services struct {
	Document document.ServiceItf
	QA       *qa.Service
}

func NewApp(cfg config.Config) (*App, error) {
	if strings.TrimSpace(cfg.OpenAIAPIKey) == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}
	if cfg.DocumentConfig == nil || cfg.DocumentConfig.DocumentChunkSize <= 0 {
		return nil, fmt.Errorf("DOCUMENT_CHUNK_SIZE must be greater than zero")
	}
	if cfg.DocumentConfig.DocumentChunkOverlap < 0 || cfg.DocumentConfig.DocumentChunkOverlap >= cfg.DocumentConfig.DocumentChunkSize {
		return nil, fmt.Errorf("DOCUMENT_CHUNK_OVERLAP must be at least zero and less than DOCUMENT_CHUNK_SIZE")
	}

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
	services := InitServices(repos, &cfg)

	// Initialize the HTTP serve mux.
	mux := InitHandlers(services)

	return &App{
		DB:       db,
		Repos:    repos,
		Services: services,
		QA:       services.QA,
		Mux:      mux,
	}, nil
}

func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Document: document.NewRepository(db),
	}
}

func InitServices(repos *Repositories, cfg *config.Config) *Services {
	embedder := embedding.NewOpenAIClient(cfg.OpenAIAPIKey)
	return &Services{
		Document: document.NewService(
			repos.Document,
			cfg.DocumentConfig,
			embedder,
		),
		QA: qa.NewService(
			repos.Document,
			embedder,
			qa.NewOpenAIAnswerer(cfg.OpenAIAPIKey, cfg.OpenAIChatModel),
		),
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
	handlerQA := qa.NewHandler(services.QA)

	// Initialize handler paths.
	// TODO: move route registration to feature handlers as routes grow.
	mux.HandleFunc("POST /documents", handlerDocument.CreateDocument)
	mux.HandleFunc("GET /documents", handlerDocument.ListDocuments)
	mux.HandleFunc("GET /documents/{id}", handlerDocument.GetDocument)
	mux.HandleFunc("POST /documents/{id}/reindex", handlerDocument.ReindexDocument)
	mux.HandleFunc("POST /ask", handlerQA.Ask)

	return mux
}
