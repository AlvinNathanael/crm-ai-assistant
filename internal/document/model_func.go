package document

// Repositories

type RepoCreateDocumentParam struct {
	Document Document
}

type RepoCreateDocumentResult struct {
	ID int64
}

type RepoListDocumentParam struct {
	ID int64
}

type RepoListDocumentResult struct {
	Documents []Document
}

// Services

type SrvCreateDocumentParam struct {
	Document Document
}

type SrvCreateDocumentResult struct {
	ID int64
}

type SrvListDocumentParam struct {
	ID int64
}

type SrvListDocumentResult struct {
	Documents []Document
}
