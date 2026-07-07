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

const createDocumentQuery = `
	INSERT INTO documents (title, content)
	VALUES ($1, $2)
	RETURNING id;
`

func (r *Repository) CreateDocument(ctx context.Context, param RepoCreateDocumentParam) (*RepoCreateDocumentResult, error) {
	var newID int64
	err := r.db.QueryRowContext(
		ctx, createDocumentQuery, param.Title, param.Content,
	).Scan(&newID)
	if err != nil {
		log.Println("Error creating document:", err)
		return nil, err
	}
	return &RepoCreateDocumentResult{
		ID: newID,
	}, nil
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
