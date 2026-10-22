# CRM AI Assistant

CRM AI Assistant is a document-grounded question-answering service for customer relationship management teams. It indexes internal guidance and answers questions using relevant passages from the source documents, with citations that make answers easier to verify.

## Capabilities

- Ingest, retrieve, and list business documents through a REST API.
- Split document text into configurable, paragraph-aware chunks.
- Generate embeddings with OpenAI and store them in PostgreSQL with pgvector.
- Retrieve relevant passages and answer questions with source references.
- Reindex saved documents when chunking settings change.
- Run the API and database locally with Docker and Go.

## Requirements

- Go 1.24.8 or newer
- Docker Desktop with its Linux engine running, or Docker Engine
- An OpenAI API key with access to the configured embedding and response models

## Configuration

Create a local environment file from the example:

```powershell
Copy-Item .env.example .env
```

Set `OPENAI_API_KEY` in `.env`. The file also contains the PostgreSQL connection settings, the answer model, and document chunking defaults. Keep `.env` private and do not commit it.

The embedding model is `text-embedding-3-small` (1,536 dimensions). The answer model defaults to `gpt-6-luna`; change it with `OPENAI_CHAT_MODEL` if needed. Chunk size and overlap are controlled by `DOCUMENT_CHUNK_SIZE` and `DOCUMENT_CHUNK_OVERLAP`.

## Run locally on Windows

Start Docker Desktop and wait until its Linux engine is ready. From the repository root, run:

```powershell
.\run.ps1
```

The script loads `.env`, starts PostgreSQL, waits for database readiness, applies migrations, and starts the API at `http://localhost:8080`. Keep the terminal open while using the service.

In a second PowerShell terminal, run the end-to-end smoke check:

```powershell
.\scripts\smoke-test.ps1
```

The smoke check creates a sample document, asks a question, and verifies that the new document appears in the answer sources. It requires a working OpenAI API key.

## API and database

See [API Reference](API.md) for endpoint details, request and response examples, error behavior, and commands for inspecting PostgreSQL data.

The API supports health checks, document creation and retrieval, document reindexing, and question answering. Document creation stores the document, its chunks, and their embeddings in a database transaction. Reindexing replaces a document's chunks atomically; if reindexing fails, its existing chunks are retained.

Database schema changes are versioned in `migrations/`. The local PostgreSQL service uses the `pgvector/pgvector:pg17` image.

## Repository layout

```text
cmd/api/             HTTP API entrypoint
docker/              Local PostgreSQL Compose configuration
internal/app/        Application wiring and routes
internal/config/     Environment configuration
internal/database/   PostgreSQL connection setup
internal/document/   Document API, chunking, persistence, retrieval
internal/embedding/  OpenAI embeddings client
internal/qa/         Retrieval and question answering
migrations/          Database schema migrations
scripts/             Local workflow helpers
```

## Project status

The document ingestion, vector indexing, retrieval, and question-answering flows are implemented. The database migration has been applied locally, and the API health and document-list endpoints have been checked against the local database. The remaining release check is a live end-to-end run with a valid OpenAI API key, including embedding a document and receiving an answer with a source reference.

## License

No license has been declared yet. Add a `LICENSE` file before distributing this project publicly.
