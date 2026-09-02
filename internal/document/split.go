package document

import "strings"

func Split(content string, maxChunkSize int) []Chunk {
	chunks := make([]Chunk, 0)

	if maxChunkSize <= 0 {
		return chunks
	}

	// Handles Unicode whitespace, tabs, and newlines.
	normalized := strings.Join(strings.Fields(content), " ")
	if normalized == "" {
		return chunks
	}

	remaining := []rune(normalized)

	for len(remaining) > 0 {
		if len(remaining) <= maxChunkSize {
			chunks = append(chunks, Chunk{
				Index:   len(chunks),
				Content: string(remaining),
			})
			break
		}

		// Best case: a word ends exactly at the maximum size.
		// Example: "The quick brown fox", max 15.
		if remaining[maxChunkSize] == ' ' {
			chunks = append(chunks, Chunk{
				Index:   len(chunks),
				Content: string(remaining[:maxChunkSize]),
			})

			remaining = remaining[maxChunkSize+1:]
			continue
		}

		// Otherwise, find the last whitespace before the limit.
		whitespaceIndex := -1
		for i := maxChunkSize - 1; i >= 0; i-- {
			if remaining[i] == ' ' {
				whitespaceIndex = i
				break
			}
		}

		if whitespaceIndex != -1 {
			chunks = append(chunks, Chunk{
				Index:   len(chunks),
				Content: string(remaining[:whitespaceIndex]),
			})

			remaining = remaining[whitespaceIndex+1:]
			continue
		}

		// No whitespace is available: the current word is too long.
		chunks = append(chunks, Chunk{
			Index:   len(chunks),
			Content: string(remaining[:maxChunkSize]),
		})

		remaining = remaining[maxChunkSize:]
	}

	return chunks
}
