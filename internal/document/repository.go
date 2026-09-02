package document

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

var ErrDocumentNotFound = errors.New("document not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

const (
	createDocumentQuery = `
	INSERT INTO documents (title, content)
	VALUES ($1, $2)
	RETURNING id;
`

	createDocumentChunksQuery = `
	INSERT INTO document_chunks (document_id, chunk_index, content)
	VALUES ($1, $2, $3);
`
)

func (r *Repository) CreateDocumentWithChunks(ctx context.Context, param RepoCreateDocumentWithChunksParam) (*RepoCreateDocumentWithChunksResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning transaction for creating document with chunks: %v\n", err)
		return nil, err
	}
	defer tx.Rollback() // Ensure rollback in case of error.

	var newID int64
	err = tx.QueryRowContext(ctx, createDocumentQuery, param.Title, param.Content).Scan(&newID)
	if err != nil {
		log.Printf("Error creating document with chunks: %v\n", err)
		return nil, err
	}

	for _, chunk := range param.Chunks {
		_, err = tx.ExecContext(ctx, createDocumentChunksQuery, newID, chunk.Index, chunk.Content)
		if err != nil {
			log.Printf("Error creating chunk for document %d, chunk %d: %v\n", newID, chunk.Index, err)
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction for creating document %d with chunks: %v\n", newID, err)
		return nil, err
	}
	return &RepoCreateDocumentWithChunksResult{ID: newID}, nil
}

const listDocumentsQuery = `
	SELECT id, title, content
	FROM documents
	ORDER BY id DESC;
`

func (r *Repository) ListDocuments(ctx context.Context, param RepoListDocumentsParam) (*RepoListDocumentsResult, error) {
	rows, err := r.db.QueryContext(
		ctx, listDocumentsQuery,
	)
	if err != nil {
		log.Println("Error listing documents:", err)
		return nil, err
	}
	defer rows.Close()

	documents := make([]Document, 0)
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Content); err != nil {
			log.Println("Error scanning document:", err)
			return nil, err
		}
		documents = append(documents, doc)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return &RepoListDocumentsResult{
		Documents: documents,
	}, nil
}

const getDocumentQuery = `
	SELECT id, title, content
	FROM documents
	WHERE id = $1;
`

func (r *Repository) GetDocument(ctx context.Context, param RepoGetDocumentParam) (*RepoGetDocumentResult, error) {
	var doc Document
	err := r.db.QueryRowContext(
		ctx, getDocumentQuery, param.ID,
	).Scan(&doc.ID, &doc.Title, &doc.Content)
	if err != nil {
		log.Println("Error getting document:", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDocumentNotFound // document not found
		}
		return nil, err
	}

	return &RepoGetDocumentResult{
		Document: &doc,
	}, nil
}
