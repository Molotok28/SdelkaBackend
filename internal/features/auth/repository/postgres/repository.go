package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
)

// ErrNotFound возвращается, когда пользователь не найден в базе.
var ErrNotFound = errors.New("user not found")

type AuthPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewAuthPostgresRepository(pool *pgxpool.Pool) *AuthPostgresRepository {
	return &AuthPostgresRepository{pool: pool}
}

// CreateUser вставляет нового пользователя с хешем пароля и возвращает созданную запись.
func (r *AuthPostgresRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	const query = `
		INSERT INTO sdelka.users (name, surname, phone_number, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, name, surname, phone_number, password_hash
	`

	var created domain.User
	err := r.pool.QueryRow(ctx, query, user.Name, user.Surname, user.PhoneNumber, user.PasswordHash).
		Scan(&created.ID, &created.Version, &created.Name, &created.Surname, &created.PhoneNumber, &created.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	return &created, nil
}

// FindUserByPhone ищет пользователя по номеру телефона. Возвращает ErrNotFound, если не найден.
func (r *AuthPostgresRepository) FindUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	const query = `
		SELECT id, version, name, surname, phone_number, password_hash
		FROM sdelka.users
		WHERE phone_number = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, phone).
		Scan(&user.ID, &user.Version, &user.Name, &user.Surname, &user.PhoneNumber, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query user by phone: %w", err)
	}

	return &user, nil
}

// FindUserByID ищет пользователя по ID. Возвращает ErrNotFound, если не найден.
func (r *AuthPostgresRepository) FindUserByID(ctx context.Context, id int) (*domain.User, error) {
	const query = `
		SELECT id, version, name, surname, phone_number, password_hash
		FROM sdelka.users
		WHERE id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Version, &user.Name, &user.Surname, &user.PhoneNumber, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}

	return &user, nil
}
