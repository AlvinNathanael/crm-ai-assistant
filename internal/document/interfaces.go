package document

import "context"

// Repositories
type RepositoryItf interface {
	ListDocuments(ctx context.Context, param RepoListDocumentsParam) (*RepoListDocumentsResult, error)
	GetDocument(ctx context.Context, param RepoGetDocumentParam) (*RepoGetDocumentResult, error)
	CreateDocumentWithChunks(ctx context.Context, param RepoCreateDocumentWithChunksParam) (*RepoCreateDocumentWithChunksResult, error)
	ReplaceDocumentChunks(ctx context.Context, documentID int64, chunks []Chunk) error
	SearchSimilarChunks(ctx context.Context, embedding []float32, limit int) ([]SearchResult, error)
}

// Services
type ServiceItf interface {
	CreateDocument(ctx context.Context, param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error)
	ListDocuments(ctx context.Context, param SrvListDocumentsParam) (*SrvListDocumentsResult, error)
	GetDocument(ctx context.Context, param SrvGetDocumentParam) (*SrvGetDocumentResult, error)
	ReindexDocument(ctx context.Context, param SrvReindexDocumentParam) (*SrvReindexDocumentResult, error)
}
