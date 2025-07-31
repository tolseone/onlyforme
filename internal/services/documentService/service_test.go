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
		// Основные сценарии для url1
		{
			name: "First document for URL1",
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
			name: "Newer version for URL1",
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
			name: "Older version for URL1",
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
			name: "Middle version for URL1 - no update",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   150,
				FetchTime: 150,
				Text:      "MiddleText",
			},
			expected: nil,
		},
		{
			name: "Duplicate FetchTime for URL1 - no update",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   200,
				FetchTime: 200,
				Text:      "Duplicate",
			},
			expected: nil,
		},

		// Сценарии для другого URL (url2)
		{
			name: "First document for URL2",
			input: domain.TDocument{
				Url:       "url2",
				PubDate:   300,
				FetchTime: 300,
				Text:      "URL2 Text1",
			},
			expected: &domain.TDocument{
				Url:            "url2",
				PubDate:        300,
				FetchTime:      300,
				Text:           "URL2 Text1",
				FirstFetchTime: 300,
			},
		},
		{
			name: "Newer version for URL2",
			input: domain.TDocument{
				Url:       "url2",
				PubDate:   400,
				FetchTime: 400,
				Text:      "URL2 Text2",
			},
			expected: &domain.TDocument{
				Url:            "url2",
				PubDate:        300,
				FetchTime:      400,
				Text:           "URL2 Text2",
				FirstFetchTime: 300,
			},
		},
		{
			name: "URL1 remains unchanged after URL2 updates",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   200,
				FetchTime: 200,
				Text:      "Duplicate",
			},
			expected: nil, // Проверяем, что URL1 не изменился
		},

		// Специальные сценарии
		{
			name: "New URL with FetchTime between existing URLs",
			input: domain.TDocument{
				Url:       "url3",
				PubDate:   250,
				FetchTime: 250,
				Text:      "URL3 Text",
			},
			expected: &domain.TDocument{
				Url:            "url3",
				PubDate:        250,
				FetchTime:      250,
				Text:           "URL3 Text",
				FirstFetchTime: 250,
			},
		},
		{
			name: "Empty text update",
			input: domain.TDocument{
				Url:       "url1",
				PubDate:   500,
				FetchTime: 500,
				Text:      "", // Пустой текст
			},
			expected: &domain.TDocument{
				Url:            "url1",
				PubDate:        50,  // Сохраняется минимальный
				FetchTime:      500, // Новый максимум
				Text:           "",  // Пустой текст
				FirstFetchTime: 50,
			},
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
