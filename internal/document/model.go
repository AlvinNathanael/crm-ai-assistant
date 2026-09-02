package document

type Document struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Chunk struct {
	Index   int
	Content string
}
