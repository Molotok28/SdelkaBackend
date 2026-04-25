package ads_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
	ads_postgres_repository "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/repository/postgres"
)

var (
	ErrNotFound  = errors.New("ad not found")
	ErrForbidden = errors.New("you don't own this ad")
)

// AdsRepository — интерфейс репозитория, используемый сервисным слоем.
type AdsRepository interface {
	Create(ctx context.Context, ad domain.Ad) (*domain.Ad, error)
	FindByID(ctx context.Context, id int) (*domain.Ad, error)
	FindAll(ctx context.Context, limit, offset int) ([]*domain.Ad, int, error)
	Update(ctx context.Context, ad domain.Ad) (*domain.Ad, error)
	Delete(ctx context.Context, id int) error
}

type CreateAdInput struct {
	Title       string
	Description string
	Price       float64
	UserID      int
}

type UpdateAdInput struct {
	ID          int
	Title       string
	Description string
	Price       float64
	RequesterID int
}

type ListAdsInput struct {
	Limit  int
	Offset int
}

type ListAdsResult struct {
	Items  []*domain.Ad
	Total  int
	Limit  int
	Offset int
}

type AdsService struct {
	repo AdsRepository
}

func NewAdsService(repo AdsRepository) *AdsService {
	return &AdsService{repo: repo}
}

// CreateAd создаёт новое объявление от имени аутентифицированного пользователя.
func (s *AdsService) CreateAd(ctx context.Context, input CreateAdInput) (*domain.Ad, error) {
	ad := domain.Ad{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		UserID:      input.UserID,
	}

	created, err := s.repo.Create(ctx, ad)
	if err != nil {
		return nil, fmt.Errorf("create ad: %w", err)
	}

	return created, nil
}

// GetAd возвращает объявление по ID. Доступно публично.
func (s *AdsService) GetAd(ctx context.Context, id int) (*domain.Ad, error) {
	ad, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ads_postgres_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get ad: %w", err)
	}

	return ad, nil
}

// ListAds возвращает страницу объявлений. Доступно публично.
func (s *AdsService) ListAds(ctx context.Context, input ListAdsInput) (ListAdsResult, error) {
	items, total, err := s.repo.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		return ListAdsResult{}, fmt.Errorf("list ads: %w", err)
	}

	return ListAdsResult{
		Items:  items,
		Total:  total,
		Limit:  input.Limit,
		Offset: input.Offset,
	}, nil
}

// UpdateAd обновляет объявление. Возвращает ErrForbidden, если requesterID не совпадает с автором.
func (s *AdsService) UpdateAd(ctx context.Context, input UpdateAdInput) (*domain.Ad, error) {
	existing, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, ads_postgres_repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find ad for update: %w", err)
	}

	if existing.UserID != input.RequesterID {
		return nil, ErrForbidden
	}

	updated, err := s.repo.Update(ctx, domain.Ad{
		ID:          input.ID,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
	})
	if err != nil {
		return nil, fmt.Errorf("update ad: %w", err)
	}

	return updated, nil
}

// DeleteAd удаляет объявление. Возвращает ErrForbidden, если requesterID не совпадает с автором.
func (s *AdsService) DeleteAd(ctx context.Context, id, requesterID int) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ads_postgres_repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("find ad for delete: %w", err)
	}

	if existing.UserID != requesterID {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete ad: %w", err)
	}

	return nil
}
