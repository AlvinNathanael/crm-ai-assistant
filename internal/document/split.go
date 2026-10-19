package document

import "strings"

// Split divides content into chunks of at most maxChunkSize runes. Adjacent
// chunks repeat up to overlap runes of context. Invalid sizes or overlap
// values return an empty slice.
func Split(content string, maxChunkSize, overlap int) []Chunk {
	chunks := make([]Chunk, 0)

	if maxChunkSize <= 0 || overlap < 0 || overlap >= maxChunkSize {
		return chunks
	}

	// Handles Unicode whitespace, tabs, and newlines.
	normalized := strings.Join(strings.Fields(content), " ")
	if normalized == "" {
		return chunks
	}

	text := []rune(normalized)
	start := 0

	for start < len(text) {
		remaining := text[start:]
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

			start = nextChunkStart(text, start, start+maxChunkSize+1, overlap)
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

			start = nextChunkStart(text, start, start+whitespaceIndex+1, overlap)
			continue
		}

		// No whitespace is available: the current word is too long.
		chunks = append(chunks, Chunk{
			Index:   len(chunks),
			Content: string(remaining[:maxChunkSize]),
		})

		start = nextChunkStart(text, start, start+maxChunkSize, overlap)
	}

	return chunks
}

func nextChunkStart(text []rune, currentStart, chunkEnd, overlap int) int {
	next := chunkEnd - overlap
	// A short chunk caused by a nearby word boundary must still make progress.
	if next <= currentStart {
		next = currentStart + 1
	}
	// Keep the overlap on a word boundary when possible instead of starting
	// the next chunk in the middle of a word.
	if next > currentStart && next < chunkEnd && text[next] != ' ' {
		for next > currentStart && text[next-1] != ' ' {
			next--
		}
		if next == currentStart {
			next = chunkEnd - overlap
			if next <= currentStart {
				next = currentStart + 1
			}
		}
	}
	for next < chunkEnd && text[next] == ' ' {
		next++
	}
	return next
}
