package document

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	INSERT INTO document_chunks (document_id, chunk_index, content, embedding)
	VALUES ($1, $2, $3, $4::vector);
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

	if err := insertDocumentChunks(ctx, tx, newID, param.Chunks); err != nil {
		log.Printf("Error creating chunks for document %d: %v\n", newID, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction for creating document %d with chunks: %v\n", newID, err)
		return nil, err
	}
	return &RepoCreateDocumentWithChunksResult{ID: newID}, nil
}

func insertDocumentChunks(ctx context.Context, tx *sql.Tx, documentID int64, chunks []Chunk) error {
	for _, chunk := range chunks {
		if _, err := tx.ExecContext(ctx, createDocumentChunksQuery, documentID, chunk.Index, chunk.Content, formatVector(chunk.Embedding)); err != nil {
			return fmt.Errorf("insert chunk %d: %w", chunk.Index, err)
		}
	}
	return nil
}

func (r *Repository) ReplaceDocumentChunks(ctx context.Context, documentID int64, chunks []Chunk) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM document_chunks WHERE document_id = $1`, documentID); err != nil {
		return err
	}
	if err := insertDocumentChunks(ctx, tx, documentID, chunks); err != nil {
		return err
	}
	return tx.Commit()
}

func formatVector(values []float32) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = strconv.FormatFloat(float64(value), 'g', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

const searchSimilarChunksQuery = `
	SELECT dc.document_id, d.title, dc.chunk_index, dc.content,
	       1 - (dc.embedding <=> $1::vector) AS similarity
	FROM document_chunks AS dc
	JOIN documents AS d ON d.id = dc.document_id
	WHERE dc.embedding IS NOT NULL
	ORDER BY dc.embedding <=> $1::vector
	LIMIT $2;
`

func (r *Repository) SearchSimilarChunks(ctx context.Context, embedding []float32, limit int) ([]SearchResult, error) {
	rows, err := r.db.QueryContext(ctx, searchSimilarChunksQuery, formatVector(embedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]SearchResult, 0, limit)
	for rows.Next() {
		var result SearchResult
		if err := rows.Scan(
			&result.DocumentID,
			&result.Title,
			&result.ChunkIndex,
			&result.Content,
			&result.Similarity,
		); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
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
