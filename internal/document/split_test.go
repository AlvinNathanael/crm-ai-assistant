package document_test

import (
	"reflect"
	"testing"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/document"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		maxChunkSize int
		overlap      int
		want         []document.Chunk
	}{
		{
			name:         "empty content returns empty slice",
			content:      "",
			maxChunkSize: 10,
			want:         []document.Chunk{},
		},
		{
			name:         "whitespace only content returns empty slice",
			content:      " \n\t  ",
			maxChunkSize: 10,
			want:         []document.Chunk{},
		},
		{
			name:         "short content returns one chunk",
			content:      "Hello CRM",
			maxChunkSize: 20,
			want: []document.Chunk{
				{Index: 0, Content: "Hello CRM"},
			},
		},
		{
			name:         "trims leading and trailing whitespace",
			content:      "  Hello CRM users.  ",
			maxChunkSize: 30,
			want: []document.Chunk{
				{Index: 0, Content: "Hello CRM users."},
			},
		},
		{
			name:         "splits at word boundaries and preserves order",
			content:      "The quick brown fox jumps over the lazy dog",
			maxChunkSize: 15,
			want: []document.Chunk{
				{Index: 0, Content: "The quick brown"},
				{Index: 1, Content: "fox jumps over"},
				{Index: 2, Content: "the lazy dog"},
			},
		},
		{
			name:         "does not leave leading or trailing whitespace in chunks",
			content:      "  The quick brown fox  ",
			maxChunkSize: 10,
			want: []document.Chunk{
				{Index: 0, Content: "The quick"},
				{Index: 1, Content: "brown fox"},
			},
		},
		{
			name:         "normalizes repeated whitespace between words",
			content:      "Hello   CRM\n\tassistant",
			maxChunkSize: 30,
			want: []document.Chunk{
				{Index: 0, Content: "Hello CRM assistant"},
			},
		},
		{
			name:         "reflows wrapped lines and preserves paragraph boundaries",
			content:      "Lead response policy\r\nFirst attempt within one business day.\r\n\r\nOpportunity policy\r\nEvery open opportunity needs an owner.",
			maxChunkSize: 200,
			want: []document.Chunk{
				{Index: 0, Content: "Lead response policy First attempt within one business day.\n\nOpportunity policy Every open opportunity needs an owner."},
			},
		},
		{
			name:         "paragraph boundary remains in normalized chunk",
			content:      "First paragraph.\n\nSecond paragraph.",
			maxChunkSize: 100,
			want: []document.Chunk{
				{Index: 0, Content: "First paragraph.\n\nSecond paragraph."},
			},
		},
		{
			name:         "splits a word only when word exceeds maximum size",
			content:      "supercalifragilistic",
			maxChunkSize: 10,
			want: []document.Chunk{
				{Index: 0, Content: "supercalif"},
				{Index: 1, Content: "ragilistic"},
			},
		},
		{
			name:         "splits long word then continues with following words",
			content:      "supercalifragilistic CRM",
			maxChunkSize: 10,
			want: []document.Chunk{
				{Index: 0, Content: "supercalif"},
				{Index: 1, Content: "ragilistic"},
				{Index: 2, Content: "CRM"},
			},
		},
		{
			name:         "zero maximum size returns empty slice",
			content:      "Hello CRM",
			maxChunkSize: 0,
			want:         []document.Chunk{},
		},
		{
			name:         "negative maximum size returns empty slice",
			content:      "Hello CRM",
			maxChunkSize: -1,
			want:         []document.Chunk{},
		},
		{
			name:         "overlap repeats preceding word context",
			content:      "alpha beta gamma delta",
			maxChunkSize: 11,
			overlap:      4,
			want: []document.Chunk{
				{Index: 0, Content: "alpha beta"},
				{Index: 1, Content: "beta gamma"},
				{Index: 2, Content: "gamma delta"},
			},
		},
		{
			name:         "zero overlap preserves non-overlapping chunks",
			content:      "The quick brown fox",
			maxChunkSize: 10,
			overlap:      0,
			want: []document.Chunk{
				{Index: 0, Content: "The quick"},
				{Index: 1, Content: "brown fox"},
			},
		},
		{
			name:         "negative overlap returns empty slice",
			content:      "Hello CRM",
			maxChunkSize: 10,
			overlap:      -1,
			want:         []document.Chunk{},
		},
		{
			name:         "overlap equal to maximum returns empty slice",
			content:      "Hello CRM",
			maxChunkSize: 10,
			overlap:      10,
			want:         []document.Chunk{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := document.Split(tt.content, tt.maxChunkSize, tt.overlap)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
