package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
)

type UserPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserPostgresRepository(pool *pgxpool.Pool) *UserPostgresRepository {
	return &UserPostgresRepository{pool: pool}
}

// Create вставляет нового пользователя в таблицу sdelka.users и возвращает его с присвоенными id и version.
func (r *UserPostgresRepository) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	const query = `
		INSERT INTO sdelka.users (name, surname, phone_number)
		VALUES ($1, $2, $3)
		RETURNING id, version, name, surname, phone_number
	`

	var created domain.User
	err := r.pool.QueryRow(ctx, query, user.Name, user.Surname, user.PhoneNumber).
		Scan(&created.ID, &created.Version, &created.Name, &created.Surname, &created.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return &created, nil
}
