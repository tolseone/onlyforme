package documentRepository_test

import (
	"context"
	"database/sql"
	"testing"

	"qwerty/internal/domain/document"
	"qwerty/internal/repositories/documentRepository"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository(t *testing.T) {
	// Подключение к тестовой БД (в реальности используйте docker/testcontainers)
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/testdb?sslmode=disable")
	require.NoError(t, err)

	// Очистка БД перед тестом
	_, err = db.Exec("TRUNCATE TABLE documents RESTART IDENTITY")
	require.NoError(t, err)

	repo := documentRepository.NewPostgresRepository(db)

	tests := []struct {
		name        string
		input       document.Document
		expectError bool
		expected    *document.Document
	}{
		{
			name: "First document",
			input: document.Document{
				URL:         "https://example.com",
				PublishDate: 100,
				FetchTime:   100,
				Content:     "First content",
			},
			expected: &document.Document{
				URL:            "https://example.com",
				PublishDate:    100,
				FetchTime:      100,
				Content:        "First content",
				FirstFetchTime: 100,
			},
		},
		{
			name: "Newer document",
			input: document.Document{
				URL:         "https://example.com",
				PublishDate: 200,
				FetchTime:   200,
				Content:     "Updated content",
			},
			expected: &document.Document{
				URL:            "https://example.com",
				PublishDate:    100, // от первого документа
				FetchTime:      200,
				Content:        "Updated content",
				FirstFetchTime: 100,
			},
		},
		{
			name: "Older document",
			input: document.Document{
				URL:         "https://example.com",
				PublishDate: 50,
				FetchTime:   50,
				Content:     "Old content",
			},
			expected: &document.Document{
				URL:            "https://example.com",
				PublishDate:    50,  // новый минимум
				FetchTime:      200, // последний
				Content:        "Updated content",
				FirstFetchTime: 50, // новый минимум
			},
		},
		{
			name: "Duplicate document",
			input: document.Document{
				URL:         "https://example.com",
				PublishDate: 200,
				FetchTime:   200,
				Content:     "Duplicate",
			},
			expected: nil, // не должен вставиться
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, updated, err := repo.ComputeAggregated(context.Background(), tt.input)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.expected == nil {
				require.False(t, updated)
				require.Nil(t, result)
			} else {
				require.True(t, updated)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
