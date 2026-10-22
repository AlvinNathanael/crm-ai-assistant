package qa

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
)

type testEmbedder struct {
	inputs  [][]string
	vectors [][]float32
	err     error
}

func (e *testEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	e.inputs = append(e.inputs, append([]string(nil), inputs...))
	return e.vectors, e.err
}

type testRetriever struct {
	embedding []float32
	limit     int
	results   []document.SearchResult
	err       error
}

func (r *testRetriever) SearchSimilarChunks(_ context.Context, vector []float32, limit int) ([]document.SearchResult, error) {
	r.embedding = vector
	r.limit = limit
	return r.results, r.err
}

type testAnswerer struct {
	question string
	sources  []document.SearchResult
	answer   string
	err      error
	calls    int
}

func (a *testAnswerer) Answer(_ context.Context, question string, sources []document.SearchResult) (string, error) {
	a.calls++
	a.question = question
	a.sources = sources
	return a.answer, a.err
}

func TestAskEmbedsRetrievesAndReturnsSources(t *testing.T) {
	sources := []document.SearchResult{{
		DocumentID: 9,
		Title:      "Lead response policy",
		ChunkIndex: 0,
		Content:    "Contact a new lead within one business day.",
		Similarity: 0.91,
	}}
	embedder := &testEmbedder{vectors: [][]float32{{0.1, 0.2, 0.3}}}
	retriever := &testRetriever{results: sources}
	answerer := &testAnswerer{answer: "Within one business day."}
	service := NewService(retriever, embedder, answerer)

	result, err := service.Ask(context.Background(), "When should we contact a new lead?")
	if err != nil {
		t.Fatalf("Ask() error = %v", err)
	}
	if !reflect.DeepEqual(embedder.inputs, [][]string{{"When should we contact a new lead?"}}) {
		t.Errorf("embedded inputs = %#v", embedder.inputs)
	}
	if !reflect.DeepEqual(retriever.embedding, []float32{0.1, 0.2, 0.3}) {
		t.Errorf("retriever embedding = %#v", retriever.embedding)
	}
	if retriever.limit != defaultSourceLimit {
		t.Errorf("retrieval limit = %d, want %d", retriever.limit, defaultSourceLimit)
	}
	if answerer.question != "When should we contact a new lead?" || !reflect.DeepEqual(answerer.sources, sources) {
		t.Errorf("answerer received question %q and sources %#v", answerer.question, answerer.sources)
	}
	if result.Answer != "Within one business day." || !reflect.DeepEqual(result.Sources, sources) {
		t.Errorf("Ask() = %#v", result)
	}
}

func TestAskWithNoIndexedChunksSkipsAnswerGeneration(t *testing.T) {
	embedder := &testEmbedder{vectors: [][]float32{{0.1}}}
	answerer := &testAnswerer{answer: "should not be called"}
	service := NewService(&testRetriever{}, embedder, answerer)

	result, err := service.Ask(context.Background(), "A question")
	if err != nil {
		t.Fatalf("Ask() error = %v", err)
	}
	if answerer.calls != 0 {
		t.Errorf("answerer calls = %d, want 0", answerer.calls)
	}
	if result.Answer == "" || len(result.Sources) != 0 {
		t.Errorf("Ask() = %#v, want explanatory answer and no sources", result)
	}
}

func TestAskPropagatesEmbeddingFailureAsUpstreamError(t *testing.T) {
	embedder := &testEmbedder{err: errors.New("provider unavailable")}
	service := NewService(&testRetriever{}, embedder, &testAnswerer{})

	_, err := service.Ask(context.Background(), "A question")
	if !errors.Is(err, ErrUpstream) {
		t.Errorf("Ask() error = %v, want ErrUpstream", err)
	}
}
