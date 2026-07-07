package document

import "context"

// Repositories
type RepositoryItf interface {
	CreateDocument(ctx context.Context, param RepoCreateDocumentParam) (*RepoCreateDocumentResult, error)
	ListDocuments(ctx context.Context, param RepoListDocumentsParam) (*RepoListDocumentsResult, error)
	GetDocument(ctx context.Context, param RepoGetDocumentParam) (*RepoGetDocumentResult, error)
}

// Services
type ServiceItf interface {
	CreateDocument(ctx context.Context, param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error)
	ListDocuments(ctx context.Context, param SrvListDocumentsParam) (*SrvListDocumentsResult, error)
	GetDocument(ctx context.Context, param SrvGetDocumentParam) (*SrvGetDocumentResult, error)
}
