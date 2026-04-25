package ads_transport_http

import (
	"context"
	"net/http"
	"time"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
	core_http_middleware "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/middleware"
	core_http_server "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/server"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
)

// AdsService — интерфейс сервисного слоя, используемый транспортом.
type AdsService interface {
	CreateAd(ctx context.Context, input ads_service.CreateAdInput) (*domain.Ad, error)
	GetAd(ctx context.Context, id int) (*domain.Ad, error)
	ListAds(ctx context.Context, input ads_service.ListAdsInput) (ads_service.ListAdsResult, error)
	UpdateAd(ctx context.Context, input ads_service.UpdateAdInput) (*domain.Ad, error)
	DeleteAd(ctx context.Context, id, requesterID int) error
}

type AdsHTTPHandler struct {
	service   AdsService
	jwtConfig core_auth.Config
}

func NewAdsHTTPHandler(service AdsService, jwtConfig core_auth.Config) *AdsHTTPHandler {
	return &AdsHTTPHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (h *AdsHTTPHandler) Routes() []core_http_server.Route {
	authMW := core_http_middleware.Auth(h.jwtConfig)

	return []core_http_server.Route{
		{Method: http.MethodGet, Path: "/ads", Handler: h.ListAds},
		{Method: http.MethodGet, Path: "/ads/{id}", Handler: h.GetAd},
		{Method: http.MethodPost, Path: "/ads", Handler: h.CreateAd, Middleware: []func(http.Handler) http.Handler{authMW}},
		{Method: http.MethodPut, Path: "/ads/{id}", Handler: h.UpdateAd, Middleware: []func(http.Handler) http.Handler{authMW}},
		{Method: http.MethodDelete, Path: "/ads/{id}", Handler: h.DeleteAd, Middleware: []func(http.Handler) http.Handler{authMW}},
	}
}

// adResponse — структура ответа для одного объявления.
type adResponse struct {
	ID          int       `json:"id"`
	Version     int       `json:"version"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func domainToResponse(ad *domain.Ad) adResponse {
	return adResponse{
		ID:          ad.ID,
		Version:     ad.Version,
		Title:       ad.Title,
		Description: ad.Description,
		Price:       ad.Price,
		UserID:      ad.UserID,
		CreatedAt:   ad.CreatedAt,
		UpdatedAt:   ad.UpdatedAt,
	}
}
