package documentService

import (
	"context"
	"qwerty/internal/domain"
	"qwerty/internal/repositories/documentRepository"
)

type DocumentProcessor struct {
	repo documentRepository.DocumentRepository
}

func NewDocumentProcessor(repo documentRepository.DocumentRepository) *DocumentProcessor {
	return &DocumentProcessor{repo: repo}
}

func (p *DocumentProcessor) Process(ctx context.Context, doc domain.TDocument) (*domain.TDocument, error) {
	aggregatedDoc, updated, err := p.repo.Upsert(ctx, doc)
	if err != nil {
		return nil, err
	}

	if !updated { // Кейс, когда мы получили дубль
		return nil, nil
	}

	return aggregatedDoc, nil
}
