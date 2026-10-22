package document

import (
	"context"
	"fmt"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/config"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/embedding"
)

type Service struct {
	repo        RepositoryItf
	documentCfg *config.DocumentConfig
	embedder    embedding.Embedder
}

func NewService(repo RepositoryItf, documentCfg *config.DocumentConfig, embedder embedding.Embedder) *Service {
	return &Service{
		repo:        repo,
		documentCfg: documentCfg,
		embedder:    embedder,
	}
}

func (s *Service) CreateDocument(ctx context.Context, param SrvCreateDocumentParam) (*SrvCreateDocumentResult, error) {
	chunks := Split(
		param.Content,
		s.documentCfg.DocumentChunkSize,
		s.documentCfg.DocumentChunkOverlap,
	)
	chunkContents := make([]string, len(chunks))
	for i, chunk := range chunks {
		chunkContents[i] = chunk.Content
	}
	embeddings, err := s.embedder.Embed(ctx, chunkContents)
	if err != nil {
		return nil, fmt.Errorf("embed document chunks: %w", err)
	}
	if len(embeddings) != len(chunks) {
		return nil, fmt.Errorf("embedder returned %d vectors for %d chunks", len(embeddings), len(chunks))
	}
	for i := range chunks {
		chunks[i].Embedding = embeddings[i]
	}

	// Generate embeddings before opening the transaction; persist the document,
	// chunks, and vectors atomically in the repository.
	repoResult, err := s.repo.CreateDocumentWithChunks(ctx, RepoCreateDocumentWithChunksParam{
		Title:   param.Title,
		Content: param.Content,
		Chunks:  chunks,
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

func (s *Service) ReindexDocument(ctx context.Context, param SrvReindexDocumentParam) (*SrvReindexDocumentResult, error) {
	documentResult, err := s.repo.GetDocument(ctx, RepoGetDocumentParam{ID: param.ID})
	if err != nil {
		return nil, err
	}

	chunks := Split(
		documentResult.Document.Content,
		s.documentCfg.DocumentChunkSize,
		s.documentCfg.DocumentChunkOverlap,
	)
	chunkContents := make([]string, len(chunks))
	for i := range chunks {
		chunkContents[i] = chunks[i].Content
	}
	embeddings, err := s.embedder.Embed(ctx, chunkContents)
	if err != nil {
		return nil, fmt.Errorf("embed existing document chunks: %w", err)
	}
	if len(embeddings) != len(chunks) {
		return nil, fmt.Errorf("embedder returned %d vectors for %d chunks", len(embeddings), len(chunks))
	}
	for i := range chunks {
		chunks[i].Embedding = embeddings[i]
	}
	if err := s.repo.ReplaceDocumentChunks(ctx, param.ID, chunks); err != nil {
		return nil, fmt.Errorf("replace document chunks: %w", err)
	}

	return &SrvReindexDocumentResult{ID: param.ID, ChunkCount: len(chunks)}, nil
}
