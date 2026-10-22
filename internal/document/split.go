package document

import (
	"strings"
	"unicode"
)

// Split divides content into chunks of at most maxChunkSize runes. Adjacent
// chunks repeat up to overlap runes of context. Invalid sizes or overlap
// values return an empty slice.
func Split(content string, maxChunkSize, overlap int) []Chunk {
	chunks := make([]Chunk, 0)

	if maxChunkSize <= 0 || overlap < 0 || overlap >= maxChunkSize {
		return chunks
	}

	// Reflow wrapped lines within each paragraph, but keep blank-line paragraph
	// boundaries so separate sections do not become one continuous sentence.
	normalized := normalizeParagraphs(content)
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
				Content: strings.TrimSpace(string(remaining)),
			})
			break
		}

		// Best case: a word ends exactly at the maximum size.
		// Example: "The quick brown fox", max 15.
		if unicode.IsSpace(remaining[maxChunkSize]) {
			chunks = append(chunks, Chunk{
				Index:   len(chunks),
				Content: strings.TrimSpace(string(remaining[:maxChunkSize])),
			})

			start = nextChunkStart(text, start, start+maxChunkSize+1, overlap)
			continue
		}

		// Otherwise, find the last whitespace before the limit.
		whitespaceIndex := -1
		for i := maxChunkSize - 1; i >= 0; i-- {
			if unicode.IsSpace(remaining[i]) {
				whitespaceIndex = i
				break
			}
		}

		if whitespaceIndex != -1 {
			chunks = append(chunks, Chunk{
				Index:   len(chunks),
				Content: strings.TrimSpace(string(remaining[:whitespaceIndex])),
			})

			start = nextChunkStart(text, start, start+whitespaceIndex+1, overlap)
			continue
		}

		// No whitespace is available: the current word is too long.
		chunks = append(chunks, Chunk{
			Index:   len(chunks),
			Content: strings.TrimSpace(string(remaining[:maxChunkSize])),
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
	if next > currentStart && next < chunkEnd && !unicode.IsSpace(text[next]) {
		for next > currentStart && !unicode.IsSpace(text[next-1]) {
			next--
		}
		if next == currentStart {
			next = chunkEnd - overlap
			if next <= currentStart {
				next = currentStart + 1
			}
		}
	}
	for next < len(text) && unicode.IsSpace(text[next]) {
		next++
	}
	return next
}

func normalizeParagraphs(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	paragraphs := make([]string, 0)
	lines := make([]string, 0)
	flushParagraph := func() {
		if len(lines) == 0 {
			return
		}
		paragraphs = append(paragraphs, strings.Join(lines, " "))
		lines = lines[:0]
	}

	for _, line := range strings.Split(content, "\n") {
		words := strings.Fields(line)
		if len(words) == 0 {
			flushParagraph()
			continue
		}
		lines = append(lines, strings.Join(words, " "))
	}
	flushParagraph()

	return strings.Join(paragraphs, "\n\n")
}
