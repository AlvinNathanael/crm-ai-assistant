package document

type Document struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Chunk struct {
	Index     int
	Content   string
	Embedding []float32
}

type SearchResult struct {
	DocumentID int64   `json:"document_id"`
	Title      string  `json:"title"`
	ChunkIndex int     `json:"chunk_index"`
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}
