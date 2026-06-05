package document

type Repository struct {
	documents []Document
	// --- TODO: db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateDocument(param RepoCreateDocumentParam) (*RepoCreateDocumentResult, error)

func (r *Repository) ListDocument(param RepoListDocumentParam) (*RepoListDocumentResult, error)
