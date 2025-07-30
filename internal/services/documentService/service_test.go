package documentService

import (
	"context"
	"qwerty/internal/domain"
	"qwerty/internal/repositories/documentRepository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	repo := documentRepository.NewMemoryRepository()
	processor := NewDocumentProcessor(repo)

	tests := []struct {
		name     string
		input    domain.TDocument
		expected *domain.TDocument
	}{
		{
			name: "First document for URL",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   100,
				FetchTime: 100,
				Text:      "Text1",
			},
			expected: &domain.TDocument{
				Url:            "url1",
				PubDate:        100,
				FetchTime:      100,
				Text:           "Text1",
				FirstFetchTime: 100,
			},
		},
		{
			name: "Newer version",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   200,
				FetchTime: 200,
				Text:      "Text2",
			},
			expected: &domain.TDocument{
				Url:            "url1",
				PubDate:        100, // Сохраняется от первой версии
				FetchTime:      200, // Обновляется от последней
				Text:           "Text2",
				FirstFetchTime: 100,
			},
		},
		{
			name: "Older version",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   50,
				FetchTime: 50,
				Text:      "OldText",
			},
			expected: &domain.TDocument{
				Url:            "url1",
				PubDate:        50,  // Новый минимум
				FetchTime:      200, // Сохраняется последний
				Text:           "Text2",
				FirstFetchTime: 50, // Новый минимум
			},
		},
		{
			name: "Middle version - no update",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   100,
				FetchTime: 150,
				Text:      "MiddleText",
			},
			expected: nil,
		},
		{
			name: "Duplicate FetchTime - no update",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   200,
				FetchTime: 200,
				Text:      "Duplicate",
			},
			expected: nil, // Не должно быть обновления
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := processor.Process(ctx, tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
