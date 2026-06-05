package document

import (
	"context"
)

type Service struct {
	repo RepositoryItf
}

func NewService(repo RepositoryItf) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateDocument(ctx context.Context, param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error) {
	// Call the repository to create the document.
	repoResult, err := s.repo.CreateDocument(ctx, RepoCreateDocumentParam{
		Document: param.Document,
	})
	if err != nil {
		return nil, err
	}

	return &SrvCreateDocumentResult{
		ID: repoResult.ID,
	}, nil
}

func (s *Service) ListDocument(ctx context.Context, param SrvListDocumentParam) (*SrvListDocumentResult, error) {
	// Call the repository to list the documents.
	repoResult, err := s.repo.ListDocument(ctx, RepoListDocumentParam{
		ID: param.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SrvListDocumentResult{
		Documents: repoResult.Documents,
	}, nil
}
