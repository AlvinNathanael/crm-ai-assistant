package document_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/config"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
)

type serviceTestEmbedder struct {
	inputs  [][]string
	vectors [][]float32
	err     error
}

func (e *serviceTestEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	e.inputs = append(e.inputs, append([]string(nil), inputs...))
	return e.vectors, e.err
}

type serviceTestRepository struct {
	createCalls int
	createParam document.RepoCreateDocumentWithChunksParam
	createErr   error
}

func (r *serviceTestRepository) ListDocuments(context.Context, document.RepoListDocumentsParam) (*document.RepoListDocumentsResult, error) {
	return nil, errors.New("unexpected ListDocuments call")
}

func (r *serviceTestRepository) GetDocument(context.Context, document.RepoGetDocumentParam) (*document.RepoGetDocumentResult, error) {
	return nil, errors.New("unexpected GetDocument call")
}

func (r *serviceTestRepository) CreateDocumentWithChunks(_ context.Context, param document.RepoCreateDocumentWithChunksParam) (*document.RepoCreateDocumentWithChunksResult, error) {
	r.createCalls++
	r.createParam = param
	if r.createErr != nil {
		return nil, r.createErr
	}
	return &document.RepoCreateDocumentWithChunksResult{ID: 42}, nil
}

func (r *serviceTestRepository) ReplaceDocumentChunks(context.Context, int64, []document.Chunk) error {
	return errors.New("unexpected ReplaceDocumentChunks call")
}

func (r *serviceTestRepository) SearchSimilarChunks(context.Context, []float32, int) ([]document.SearchResult, error) {
	return nil, errors.New("unexpected SearchSimilarChunks call")
}

func TestCreateDocumentEmbedsChunksBeforePersisting(t *testing.T) {
	content := "Lead response policy\nContact the lead within one business day.\n\nOpportunity policy\nRecord the next step."
	wantChunks := document.Split(content, 200, 20)
	wantInputs := make([]string, len(wantChunks))
	for i := range wantChunks {
		wantInputs[i] = wantChunks[i].Content
	}
	embedder := &serviceTestEmbedder{vectors: [][]float32{{0.1, 0.2, 0.3}}}
	repo := &serviceTestRepository{}
	service := document.NewService(repo, &config.DocumentConfig{
		DocumentChunkSize:    200,
		DocumentChunkOverlap: 20,
	}, embedder)

	result, err := service.CreateDocument(context.Background(), document.SrvCreateDocumentParam{
		Title:   "CRM policies",
		Content: content,
	})
	if err != nil {
		t.Fatalf("CreateDocument() error = %v", err)
	}
	if result.ID != 42 {
		t.Errorf("document ID = %d, want 42", result.ID)
	}
	if !reflect.DeepEqual(embedder.inputs, [][]string{wantInputs}) {
		t.Errorf("embedded inputs = %#v, want %#v", embedder.inputs, [][]string{wantInputs})
	}
	if repo.createCalls != 1 {
		t.Fatalf("repository create calls = %d, want 1", repo.createCalls)
	}
	if repo.createParam.Title != "CRM policies" || repo.createParam.Content != content {
		t.Errorf("persisted document = %#v", repo.createParam)
	}
	if len(repo.createParam.Chunks) != 1 || !reflect.DeepEqual(repo.createParam.Chunks[0].Embedding, []float32{0.1, 0.2, 0.3}) {
		t.Errorf("persisted chunks do not contain the generated embedding: %#v", repo.createParam.Chunks)
	}
	if repo.createParam.Chunks[0].Content != wantInputs[0] {
		t.Errorf("persisted chunk content = %q, want %q", repo.createParam.Chunks[0].Content, wantInputs[0])
	}
}

func TestCreateDocumentDoesNotPersistWhenEmbeddingFails(t *testing.T) {
	repo := &serviceTestRepository{}
	service := document.NewService(repo, &config.DocumentConfig{
		DocumentChunkSize:    100,
		DocumentChunkOverlap: 10,
	}, &serviceTestEmbedder{err: errors.New("embedding provider unavailable")})

	if _, err := service.CreateDocument(context.Background(), document.SrvCreateDocumentParam{
		Title:   "CRM policies",
		Content: "A policy paragraph.",
	}); err == nil {
		t.Fatal("CreateDocument() error = nil, want embedding error")
	}
	if repo.createCalls != 0 {
		t.Errorf("repository create calls = %d, want 0 after embedding failure", repo.createCalls)
	}
}
