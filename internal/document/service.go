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
		Title:   param.Title,
		Content: param.Content,
	})
	if err != nil {
		return nil, err
	}

	return &SrvCreateDocumentResult{
		ID: repoResult.ID,
	}, nil
}

func (s *Service) ListDocuments(ctx context.Context, param SrvListDocumentsParam) (*SrvListDocumentsResult, error) {
	// Call the repository to list the documents.
	repoResult, err := s.repo.ListDocuments(ctx, RepoListDocumentsParam{})
	if err != nil {
		return nil, err
	}

	return &SrvListDocumentsResult{
		Documents: repoResult.Documents,
	}, nil
}

func (s *Service) GetDocument(ctx context.Context, param SrvGetDocumentParam) (*SrvGetDocumentResult, error) {
	// Call the repository to get the document.
	repoResult, err := s.repo.GetDocument(ctx, RepoGetDocumentParam{
		ID: param.ID,
	})
	if err != nil {
		return nil, err
	}

	return &SrvGetDocumentResult{
		Document: repoResult.Document,
	}, nil
}
