package documentRepository

import (
	"context"
	"qwerty/internal/domain"
)

type DocumentRepository interface {
	Upsert(ctx context.Context, doc domain.TDocument) (aggregatedDoc *domain.TDocument, updated bool, err error)
	GetByURL(ctx context.Context, url string) (*domain.TDocument, error)
}
