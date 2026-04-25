package ads_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
)

// ErrNotFound возвращается, когда объявление не найдено в базе.
var ErrNotFound = errors.New("ad not found")

type AdsPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewAdsPostgresRepository(pool *pgxpool.Pool) *AdsPostgresRepository {
	return &AdsPostgresRepository{pool: pool}
}

// Create вставляет новое объявление и возвращает его с присвоенными id, version, created_at, updated_at.
func (r *AdsPostgresRepository) Create(ctx context.Context, ad domain.Ad) (*domain.Ad, error) {
	const query = `
		INSERT INTO sdelka.ads (title, description, price, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, title, description, price, user_id, created_at, updated_at
	`

	var created domain.Ad
	err := r.pool.QueryRow(ctx, query, ad.Title, ad.Description, ad.Price, ad.UserID).
		Scan(&created.ID, &created.Version, &created.Title, &created.Description,
			&created.Price, &created.UserID, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert ad: %w", err)
	}

	return &created, nil
}

// FindByID ищет объявление по ID. Возвращает ErrNotFound, если не найдено.
func (r *AdsPostgresRepository) FindByID(ctx context.Context, id int) (*domain.Ad, error) {
	const query = `
		SELECT id, version, title, description, price, user_id, created_at, updated_at
		FROM sdelka.ads
		WHERE id = $1
	`

	var ad domain.Ad
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&ad.ID, &ad.Version, &ad.Title, &ad.Description,
			&ad.Price, &ad.UserID, &ad.CreatedAt, &ad.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query ad by id: %w", err)
	}

	return &ad, nil
}

// FindAll возвращает список объявлений (новые первыми) и общее количество.
func (r *AdsPostgresRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.Ad, int, error) {
	const countQuery = `SELECT COUNT(*) FROM sdelka.ads`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count ads: %w", err)
	}

	const query = `
		SELECT id, version, title, description, price, user_id, created_at, updated_at
		FROM sdelka.ads
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query ads: %w", err)
	}
	defer rows.Close()

	var ads []*domain.Ad
	for rows.Next() {
		var ad domain.Ad
		if err := rows.Scan(&ad.ID, &ad.Version, &ad.Title, &ad.Description,
			&ad.Price, &ad.UserID, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan ad row: %w", err)
		}
		ads = append(ads, &ad)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate ad rows: %w", err)
	}

	return ads, total, nil
}

// Update обновляет поля объявления и инкрементирует version. Возвращает ErrNotFound, если не найдено.
func (r *AdsPostgresRepository) Update(ctx context.Context, ad domain.Ad) (*domain.Ad, error) {
	const query = `
		UPDATE sdelka.ads
		SET title       = $1,
		    description = $2,
		    price       = $3,
		    version     = version + 1,
		    updated_at  = NOW()
		WHERE id = $4
		RETURNING id, version, title, description, price, user_id, created_at, updated_at
	`

	var updated domain.Ad
	err := r.pool.QueryRow(ctx, query, ad.Title, ad.Description, ad.Price, ad.ID).
		Scan(&updated.ID, &updated.Version, &updated.Title, &updated.Description,
			&updated.Price, &updated.UserID, &updated.CreatedAt, &updated.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update ad: %w", err)
	}

	return &updated, nil
}

// Delete удаляет объявление по ID. Возвращает ErrNotFound, если не найдено.
func (r *AdsPostgresRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM sdelka.ads WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete ad: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
