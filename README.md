# CRM AI Assistant

A Go-based backend for a CRM-focused AI assistant that stores and retrieves business documents for downstream RAG workflows. The application is built around a document ingestion pipeline, chunk-based storage, and a simple REST API for creating, listing, and fetching documents.

## Overview

This project is a foundation for building a CRM knowledge layer that can ingest internal documents, break them into manageable chunks, and expose them through a clean API. It uses Go for the API layer, PostgreSQL for persistence, and Docker for local infrastructure.

The current implementation focuses on:

- Creating documents with title and content
- Splitting long content into chunks using configurable chunk size and overlap
- Persisting documents and chunks in PostgreSQL
- Exposing endpoints for listing and retrieving documents
- Providing a local development setup with Docker and migration support

## Tech stack

- Go
- PostgreSQL
- Docker Compose
- OpenAI configuration support for future embeddings and LLM integration

## Project structure

```text
.
├── cmd/
│   └── api/
│       └── main.go                 # Application entrypoint
├── docker/
│   └── docker-compose.yml          # Local PostgreSQL container
├── internal/
│   ├── app/
│   ├── config/
│   ├── database/
│   └── document/
├── migrations/
│   ├── *.up.sql
│   └── *.down.sql
├── .env                            # Local environment configuration
├── go.mod
├── go.sum
├── run.sh
├── README.md
└── ...
```

## Prerequisites

Before running the project, make sure you have:

- Go 1.24.8 or newer
- Docker Desktop or Docker Engine
- Access to a PostgreSQL instance (local container is included)
- An OpenAI API key if you plan to extend the app with embeddings or LLM features

## Configuration

The application loads environment variables from `.env` using `godotenv`. At minimum, configure the database settings and any OpenAI settings you plan to use.

Example:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=crm_ai_assistant
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/crm_ai_assistant?sslmode=disable
MIGRATION_DATABASE_URL=postgres://postgres:postgres@host.docker.internal:5432/crm_ai_assistant?sslmode=disable
OPENAI_API_KEY=your_openai_api_key
DOCUMENT_CHUNK_SIZE=1000
DOCUMENT_CHUNK_OVERLAP=150
```

## Running locally

### Option 1: Use the provided script

```bash
./run.sh
```

This script will:

1. Load environment variables from `.env`
2. Start the PostgreSQL container
3. Wait for PostgreSQL to become ready
4. Run database migrations
5. Start the API server with `go run ./cmd/api`

### Option 2: Start manually

```bash
docker compose -f docker/docker-compose.yml up -d
export $(grep -v '^#' .env | xargs)
go run ./cmd/api
```

The API listens on:

```text
http://localhost:8080
```

## API documentation

The API listens on `http://localhost:8080`. See [API.md](API.md) for endpoint details, request and response examples, error responses, and PowerShell commands for local testing.

## Database and migrations

Database schema is managed through SQL migration files in the `migrations/` directory. The service currently creates two tables:

- `documents`
- `document_chunks`

This gives the project a base for chunk-based retrieval and future vector-search or semantic-search improvements.

## Development notes

The app is intentionally structured around a clean internal layer:

- `cmd/api` starts the HTTP server
- `internal/app` wires dependencies together
- `internal/document` contains the document domain, handlers, and services
- `internal/database` encapsulates database access
- `migrations` defines schema changes

This makes it easier to extend the system with CRM data sources, embeddings, vector search, or additional domain-specific APIs.

## Roadmap ideas

Future enhancements could include:

- OpenAI embeddings and semantic search
- CRM-source connectors for leads, deals, contacts, and tickets
- Vector search using pgvector or a dedicated vector database
- Authentication and authorization
- Request validation and structured error handling
- Testing coverage for API and document processing

## License

This project does not currently declare a license. If you intend to share or distribute it publicly, add an appropriate `LICENSE` file and reference it here.
