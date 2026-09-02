package document

// Repositories

type RepoCreateDocumentWithChunksParam struct {
	Title   string
	Content string
	Chunks  []Chunk
}

type RepoCreateDocumentWithChunksResult struct {
	ID int64
}

type RepoListDocumentsParam struct {
}

type RepoListDocumentsResult struct {
	Documents []Document
}

type RepoGetDocumentParam struct {
	ID int64
}

type RepoGetDocumentResult struct {
	Document *Document
}

// Services

type SrvCreateDocumentParam struct {
	Title   string
	Content string
}

type SrvCreateDocumentResult struct {
	ID int64
}

type SrvListDocumentsParam struct {
}

type SrvListDocumentsResult struct {
	Documents []Document
}

type SrvGetDocumentParam struct {
	ID int64
}

type SrvGetDocumentResult struct {
	Document *Document
}

// Handlers

type HandlerCreateDocumentParam struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type HandlerCreateDocumentResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}
