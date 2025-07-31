package documentRepository

import (
	"context"
	"qwerty/internal/domain"
	"sync"
)

type minFetchTimesAndPubDates struct {
	FetchTime uint64
	PubDate   uint64
}

type MemoryRepository struct {
	mu           sync.RWMutex
	minFTsAndPDs map[string]minFetchTimesAndPubDates
	latestDocs   map[string]domain.TDocument
}

func NewMemoryRepository() DocumentRepository {
	return &MemoryRepository{
		minFTsAndPDs: make(map[string]minFetchTimesAndPubDates),
		latestDocs:   make(map[string]domain.TDocument),
	}
}

func (m *MemoryRepository) ComputeAggregated(ctx context.Context, doc domain.TDocument) (*domain.TDocument, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	url := doc.Url
	currentLatestDoc, exists := m.latestDocs[url]

	// Первый документ для URL
	if !exists {
		m.minFTsAndPDs[url] = minFetchTimesAndPubDates{
			FetchTime: doc.FetchTime,
			PubDate:   doc.PubDate,
		}
		m.latestDocs[url] = doc

		return &domain.TDocument{
			Url:            doc.Url,
			PubDate:        doc.PubDate,
			FetchTime:      doc.FetchTime,
			Text:           doc.Text,
			FirstFetchTime: doc.FetchTime,
		}, true, nil
	}

	updated := false

	// Проверка на новый минимум
	if doc.FetchTime < m.minFTsAndPDs[url].FetchTime {
		m.minFTsAndPDs[url] = minFetchTimesAndPubDates{
			FetchTime: doc.FetchTime,
			PubDate:   doc.PubDate,
		}
		updated = true
	}

	// Проверка на новый максимум
	if doc.FetchTime > currentLatestDoc.FetchTime {
		m.latestDocs[url] = doc
		updated = true
	}

	// Если не было изменений
	if !updated {
		return nil, false, nil
	}

	// Формируем агрегированный документ
	return &domain.TDocument{
		Url:            url,
		PubDate:        m.minFTsAndPDs[url].PubDate,
		FetchTime:      m.latestDocs[url].FetchTime,
		Text:           m.latestDocs[url].Text,
		FirstFetchTime: m.minFTsAndPDs[url].FetchTime,
	}, true, nil
}

func (m *MemoryRepository) GetByURL(ctx context.Context, url string) (*domain.TDocument, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	minData, exists := m.minFTsAndPDs[url]
	if !exists {
		return nil, nil
	}

	latestDoc := m.latestDocs[url] // Гарантировано существует если есть minData
	return &domain.TDocument{
		Url:            url,
		PubDate:        minData.PubDate,
		FetchTime:      latestDoc.FetchTime,
		Text:           latestDoc.Text,
		FirstFetchTime: minData.FetchTime,
	}, nil
}
