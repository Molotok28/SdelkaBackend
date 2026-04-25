package auth_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
	core_http_middleware "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/middleware"
	core_http_server "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/server"
	auth_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/service"
)

// AuthService — интерфейс сервисного слоя, необходимый транспорту.
type AuthService interface {
	Register(ctx context.Context, input auth_service.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input auth_service.LoginInput) (string, error)
	GetUser(ctx context.Context, userID int) (*domain.User, error)
}

type AuthHTTPHandler struct {
	service   AuthService
	jwtConfig core_auth.Config
}

func NewAuthHTTPHandler(service AuthService, jwtConfig core_auth.Config) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		service:   service,
		jwtConfig: jwtConfig,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	authMW := core_http_middleware.Auth(h.jwtConfig)

	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/auth/register", Handler: h.Register},
		{Method: http.MethodPost, Path: "/auth/login", Handler: h.Login},
		{Method: http.MethodPost, Path: "/auth/logout", Handler: h.Logout, Middleware: []func(http.Handler) http.Handler{authMW}},
		{Method: http.MethodGet, Path: "/auth/me", Handler: h.Me, Middleware: []func(http.Handler) http.Handler{authMW}},
	}
}
