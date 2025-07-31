package documentRepository

import (
	"context"
	"qwerty/internal/domain"
	"sync"
)

// documentGroupState - состояние группы документов
type documentGroupState struct {
	firstFetchTime uint64
	publishDate    uint64
	latestDocument domain.Document
}

// MemoryRepository - уточнение типа хранилища
type MemoryRepositoryу struct {
	mu     sync.RWMutex
	groups map[string]*documentGroupState // "groups" вместо "state"
}

func NewMemoryRepositoryу() *MemoryRepositoryу {
	return &MemoryRepositoryу{
		groups: make(map[string]*documentGroupState),
	}
}

func (m *MemoryRepositoryу) ComputeAggregated(
	ctx context.Context,
	doc domain.Document,
) (*domain.Document, bool, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	group, exists := m.groups[doc.URL]

	// Первый документ для URL
	if !exists {
		m.groups[doc.URL] = &documentGroupState{
			firstFetchTime: doc.FetchTime,
			publishDate:    doc.PubDate,
			latestDocument: doc,
		}

		return &domain.Document{
			URL:            doc.URL,
			PubDate:        doc.PubDate,
			FetchTime:      doc.FetchTime,
			Text:           doc.Text,
			FirstFetchTime: doc.FetchTime,
		}, true, nil
	}

	updated := false

	// Обновление самой ранней версии
	if doc.FetchTime < group.firstFetchTime {
		group.firstFetchTime = doc.FetchTime
		group.publishDate = doc.PubDate
		updated = true
	}

	// Обновление последней версии
	if doc.FetchTime > group.latestDocument.FetchTime {
		group.latestDocument = doc
		updated = true
	}

	if !updated {
		return nil, false, nil
	}

	// Формируем агрегированный документ
	return &domain.Document{
		URL:            doc.URL,
		PubDate:        group.publishDate,
		FetchTime:      group.latestDocument.FetchTime,
		Text:           group.latestDocument.Text,
		FirstFetchTime: group.firstFetchTime,
	}, true, nil
}
