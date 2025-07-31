package documentRepository

import (
	"context"
	"database/sql"
	"fmt"
	"qwerty/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (p *PostgresRepository) ComputeAggregated(
	ctx context.Context,
	doc domain.TDocument,
) (*domain.TDocument, bool, error) {
	// Пытаемся вставить документ
	inserted, err := p.insertDocument(ctx, doc)
	if err != nil {
		return nil, false, err
	}

	// Если дубликат - пропускаем
	if !inserted {
		return nil, false, nil
	}

	// Получаем агрегированные данные
	return p.getAggregatedDocument(ctx, doc.Url)
}

func (p *PostgresRepository) insertDocument(
	ctx context.Context,
	doc domain.TDocument,
) (bool, error) {
	query := `
        INSERT INTO documents (url, pub_date, fetch_time, text)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (url, fetch_time) DO NOTHING
    `

	result, err := p.db.ExecContext(
		ctx,
		query,
		doc.Url,
		doc.PubDate,
		doc.FetchTime,
		doc.Text,
	)
	if err != nil {
		return false, fmt.Errorf("failed to insert document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}

func (p *PostgresRepository) getAggregatedDocument(
	ctx context.Context,
	url string,
) (*domain.TDocument, bool, error) {
	query := `
        (SELECT * FROM documents WHERE url = $1 ORDER BY fetch_time ASC LIMIT 1)
        UNION ALL
        (SELECT * FROM documents WHERE url = $1 ORDER BY fetch_time DESC LIMIT 1)
    `

	rows, err := p.db.QueryContext(ctx, query, url)
	if err != nil {
		return nil, false, fmt.Errorf("failed to query aggregated data: %w", err)
	}
	defer rows.Close()

	var (
		oldestDoc  domain.TDocument
		newestDoc  domain.TDocument
		foundCount int
	)

	// Обрабатываем первую строку (самый старый документ)
	if rows.Next() {
		if err := rows.Scan(
			&oldestDoc.Url,
			&oldestDoc.PubDate,
			&oldestDoc.FetchTime,
			&oldestDoc.Text,
		); err != nil {
			return nil, false, fmt.Errorf("failed to scan oldest doc: %w", err)
		}
		foundCount++
	}

	// Обрабатываем вторую строку (самый новый документ)
	if rows.Next() {
		if err := rows.Scan(
			&newestDoc.Url,
			&newestDoc.PubDate,
			&newestDoc.FetchTime,
			&newestDoc.Text,
		); err != nil {
			return nil, false, fmt.Errorf("failed to scan newest doc: %w", err)
		}
		foundCount++
	}

	// Если нет документов для этого URL
	if foundCount == 0 {
		return nil, false, nil
	}

	// Если только одна запись (значит она и самая старая и самая новая)
	if foundCount == 1 {
		newestDoc = oldestDoc
	}

	return &domain.TDocument{
		Url:            url,
		PubDate:        oldestDoc.PubDate,   // от самой старой версии
		FetchTime:      newestDoc.FetchTime, // от самой новой версии
		Text:           newestDoc.Text,      // от самой новой версии
		FirstFetchTime: oldestDoc.FetchTime, // от самой старой версии
	}, true, nil
}
