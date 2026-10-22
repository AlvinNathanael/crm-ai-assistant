package qa

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/embedding"
)

const defaultSourceLimit = 5

var ErrUpstream = errors.New("AI provider request failed")

type Retriever interface {
	SearchSimilarChunks(ctx context.Context, embedding []float32, limit int) ([]document.SearchResult, error)
}

type Answerer interface {
	Answer(ctx context.Context, question string, sources []document.SearchResult) (string, error)
}

type Service struct {
	retriever Retriever
	embedder  embedding.Embedder
	answerer  Answerer
}

type Result struct {
	Answer  string                  `json:"answer"`
	Sources []document.SearchResult `json:"sources"`
}

func NewService(retriever Retriever, embedder embedding.Embedder, answerer Answerer) *Service {
	return &Service{retriever: retriever, embedder: embedder, answerer: answerer}
}

func (s *Service) Ask(ctx context.Context, question string) (*Result, error) {
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("question is required")
	}

	vectors, err := s.embedder.Embed(ctx, []string{question})
	if err != nil {
		return nil, fmt.Errorf("%w: create query embedding", ErrUpstream)
	}
	if len(vectors) != 1 {
		return nil, fmt.Errorf("%w: embedding provider returned an unexpected number of vectors", ErrUpstream)
	}

	sources, err := s.retriever.SearchSimilarChunks(ctx, vectors[0], defaultSourceLimit)
	if err != nil {
		return nil, fmt.Errorf("search document chunks: %w", err)
	}
	if len(sources) == 0 {
		return &Result{
			Answer:  "I couldn't find any indexed document chunks to answer from.",
			Sources: []document.SearchResult{},
		}, nil
	}

	answer, err := s.answerer.Answer(ctx, question, sources)
	if err != nil {
		return nil, fmt.Errorf("%w: generate answer", ErrUpstream)
	}
	return &Result{Answer: answer, Sources: sources}, nil
}
