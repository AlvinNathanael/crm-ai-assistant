package document

// Repositories
type RepositoryItf interface {
	CreateDocument(param RepoCreateDocumentParam) (*RepoCreateDocumentResult, error)
	ListDocument(param RepoListDocumentParam) (*RepoListDocumentResult, error)
}

// Services
type ServiceItf interface {
	CreateDocument(param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error)
	ListDocument(param SrvListDocumentParam) (*SrvListDocumentResult, error)
}
