package document

import "context"

// Repositories
type RepositoryItf interface {
	CreateDocument(ctx context.Context, param RepoCreateDocumentParam) (*RepoCreateDocumentResult, error)
	ListDocument(ctx context.Context, param RepoListDocumentParam) (*RepoListDocumentResult, error)
}

// Services
type ServiceItf interface {
	CreateDocument(ctx context.Context, param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error)
	ListDocument(ctx context.Context, param SrvListDocumentParam) (*SrvListDocumentResult, error)
}
