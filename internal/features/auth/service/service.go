package auth_service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
	auth_postgres_repository "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/repository/postgres"
)

var (
	ErrPhoneAlreadyExists = errors.New("phone number already registered")
	ErrInvalidCredentials = errors.New("invalid phone number or password")
)

// AuthRepository — интерфейс репозитория, необходимый сервису аутентификации.
type AuthRepository interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	FindUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	FindUserByID(ctx context.Context, id int) (*domain.User, error)
}

type RegisterInput struct {
	Name        string
	Surname     string
	PhoneNumber string
	Password    string
}

type LoginInput struct {
	PhoneNumber string
	Password    string
}

type AuthService struct {
	repo      AuthRepository
	jwtConfig core_auth.Config
}

func NewAuthService(repo AuthRepository, jwtConfig core_auth.Config) *AuthService {
	return &AuthService{repo: repo, jwtConfig: jwtConfig}
}

// Register регистрирует нового пользователя: проверяет уникальность телефона, хеширует пароль.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	_, err := s.repo.FindUserByPhone(ctx, input.PhoneNumber)
	if err == nil {
		return nil, ErrPhoneAlreadyExists
	}
	if !errors.Is(err, auth_postgres_repository.ErrNotFound) {
		return nil, fmt.Errorf("check phone uniqueness: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := domain.User{
		Name:         input.Name,
		Surname:      input.Surname,
		PhoneNumber:  &input.PhoneNumber,
		PasswordHash: string(hash),
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

// Login проверяет учётные данные и возвращает подписанный JWT-токен.
func (s *AuthService) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.repo.FindUserByPhone(ctx, input.PhoneNumber)
	if err != nil {
		if errors.Is(err, auth_postgres_repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := core_auth.GenerateToken(user.ID, s.jwtConfig)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

// GetUser возвращает пользователя по его ID (для эндпоинта /me).
func (s *AuthService) GetUser(ctx context.Context, userID int) (*domain.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}
