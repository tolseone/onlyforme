package documentRepository

import (
	"context"
	"qwerty/internal/domain"
)

type DocumentRepository interface {
	// ComputeAggregated - точнее отражает суть операции
	ComputeAggregated(ctx context.Context, doc domain.Document) (aggregatedDoc *domain.Document, updated bool, err error)
}
