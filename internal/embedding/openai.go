package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	model      = "text-embedding-3-small"
	dimensions = 1536
	batchSize  = 100
	endpoint   = "https://api.openai.com/v1/embeddings"
)

type Embedder interface {
	Embed(ctx context.Context, inputs []string) ([][]float32, error)
}

type OpenAIClient struct {
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:   apiKey,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (c *OpenAIClient) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	if len(inputs) == 0 {
		return [][]float32{}, nil
	}
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required to generate embeddings")
	}

	allEmbeddings := make([][]float32, len(inputs))
	for start := 0; start < len(inputs); start += batchSize {
		end := min(start+batchSize, len(inputs))
		batch, err := c.embedBatch(ctx, inputs[start:end])
		if err != nil {
			return nil, fmt.Errorf("embed input batch starting at %d: %w", start, err)
		}
		copy(allEmbeddings[start:end], batch)
	}
	return allEmbeddings, nil
}

func (c *OpenAIClient) embedBatch(ctx context.Context, inputs []string) ([][]float32, error) {
	requestBody, err := json.Marshal(struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}{
		Model: model,
		Input: inputs,
	})
	if err != nil {
		return nil, fmt.Errorf("encode embedding request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("create embedding request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send embedding request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return nil, fmt.Errorf("embedding request returned HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	var result struct {
		Data []struct {
			Index     int       `json:"index"`
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode embedding response: %w", err)
	}
	if len(result.Data) != len(inputs) {
		return nil, fmt.Errorf("embedding response returned %d vectors for %d inputs", len(result.Data), len(inputs))
	}

	embeddings := make([][]float32, len(inputs))
	for _, item := range result.Data {
		if item.Index < 0 || item.Index >= len(inputs) {
			return nil, fmt.Errorf("embedding response contains invalid input index %d", item.Index)
		}
		if len(item.Embedding) != dimensions {
			return nil, fmt.Errorf("embedding at index %d has %d dimensions; expected %d", item.Index, len(item.Embedding), dimensions)
		}
		embeddings[item.Index] = item.Embedding
	}
	for i, vector := range embeddings {
		if vector == nil {
			return nil, fmt.Errorf("embedding response is missing input index %d", i)
		}
	}

	return embeddings, nil
}
