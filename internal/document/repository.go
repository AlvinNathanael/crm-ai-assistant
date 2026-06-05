package document

import (
	"context"
)

type Repository struct {
	documents []Document // In a real application, this would be replaced with a database connection, such as *sql.DB.
	// --- TODO: db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateDocument(ctx context.Context, param RepoCreateDocumentParam) (*RepoCreateDocumentResult, error) {
	// In a real application, this would involve inserting the document into a database and returning the generated ID. For this example, we'll just append it to the in-memory slice and return a mock ID.
	newID := int64(len(r.documents) + 1)
	param.Document.ID = newID
	r.documents = append(r.documents, param.Document)
	return &RepoCreateDocumentResult{
		ID: newID,
	}, nil
}

func (r *Repository) ListDocument(ctx context.Context, param RepoListDocumentParam) (*RepoListDocumentResult, error) {
	if param.ID != 0 {
		if len(r.documents) < int(param.ID) {
			return &RepoListDocumentResult{}, nil // In a real application, you might want to return an error if the document with the specified ID does not exist.
		}
		return &RepoListDocumentResult{
			Documents: []Document{
				r.documents[param.ID-1], // This is just a mock implementation. In a real application, you would query the database for the document with the specified ID.
			},
		}, nil
	}
	return &RepoListDocumentResult{
		Documents: r.documents,
	}, nil
}
