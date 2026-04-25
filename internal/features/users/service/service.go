package users_service

import (
	"context"
	"fmt"

	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
)

// UsersRepository — интерфейс репозитория пользователей, используемый сервисным слоем.
type UsersRepository interface {
	Create(ctx context.Context, user domain.User) (*domain.User, error)
}

// CreateUserInput — входные данные для создания пользователя.
type CreateUserInput struct {
	Name        string
	Surname     string
	PhoneNumber *string
}

type UserService struct {
	repo UsersRepository
}

func NewUserService(repo UsersRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser создаёт нового пользователя и возвращает его с присвоенным ID.
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	user := domain.User{
		Name:        input.Name,
		Surname:     input.Surname,
		PhoneNumber: input.PhoneNumber,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user in repository: %w", err)
	}

	return created, nil
}
