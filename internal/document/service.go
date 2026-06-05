package document

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateDocument(param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error)

func (s *Service) ListDocument(param SrvListDocumentParam) (*SrvListDocumentResult, error)
