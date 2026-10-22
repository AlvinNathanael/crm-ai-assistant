package embedding

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestEmbedBatchesInputsAndRestoresResponseOrder(t *testing.T) {
	var batchSizes []int
	globalOffset := 0
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("request method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("authorization = %q", got)
		}
		var request struct {
			Model string   `json:"model"`
			Input []string `json:"input"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
			return nil, err
		}
		if request.Model != model {
			t.Errorf("model = %q, want %q", request.Model, model)
		}
		batchSizes = append(batchSizes, len(request.Input))
		data := make([]struct {
			Index     int       `json:"index"`
			Embedding []float32 `json:"embedding"`
		}, 0, len(request.Input))
		for i := len(request.Input) - 1; i >= 0; i-- {
			vector := make([]float32, dimensions)
			vector[0] = float32(globalOffset + i)
			data = append(data, struct {
				Index     int       `json:"index"`
				Embedding []float32 `json:"embedding"`
			}{Index: i, Embedding: vector})
		}
		globalOffset += len(request.Input)
		body, err := json.Marshal(struct {
			Data any `json:"data"`
		}{Data: data})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	})

	client := &OpenAIClient{apiKey: "test-key", endpoint: "https://embeddings.test", httpClient: &http.Client{Transport: transport}}
	inputs := make([]string, batchSize+1)
	for i := range inputs {
		inputs[i] = fmt.Sprintf("input %d", i)
	}

	got, err := client.Embed(context.Background(), inputs)
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}
	if !reflect.DeepEqual(batchSizes, []int{batchSize, 1}) {
		t.Errorf("batch sizes = %v, want [%d 1]", batchSizes, batchSize)
	}
	if len(got) != len(inputs) {
		t.Fatalf("got %d vectors, want %d", len(got), len(inputs))
	}
	for i, vector := range got {
		if len(vector) != dimensions {
			t.Fatalf("vector %d has %d dimensions, want %d", i, len(vector), dimensions)
		}
		if vector[0] != float32(i) {
			t.Errorf("vector %d first value = %v, want %d", i, vector[0], i)
		}
	}
}

func TestEmbedRejectsUnexpectedVectorDimensions(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body, err := json.Marshal(struct {
			Data []struct {
				Index     int       `json:"index"`
				Embedding []float32 `json:"embedding"`
			} `json:"data"`
		}{Data: []struct {
			Index     int       `json:"index"`
			Embedding []float32 `json:"embedding"`
		}{
			{Index: 0, Embedding: []float32{0.1, 0.2}},
		}})
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(string(body))),
			Request:    r,
		}, nil
	})

	client := &OpenAIClient{apiKey: "test-key", endpoint: "https://embeddings.test", httpClient: &http.Client{Transport: transport}}
	if _, err := client.Embed(context.Background(), []string{"one input"}); err == nil {
		t.Fatal("Embed() error = nil, want dimension validation error")
	}
}
