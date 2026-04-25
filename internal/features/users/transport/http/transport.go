package users_transport_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nikolaimoiseev01/SdelkaBackend/internal/core/domain"
	core_http_server "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/server"
	users_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/users/service"
)

type UsersHTTPHandler struct {
	userService UsersService
}

// UsersService — интерфейс сервисного слоя пользователей, используемый транспортным слоем.
type UsersService interface {
	CreateUser(ctx context.Context, input users_service.CreateUserInput) (*domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userService: usersService,
	}
}

// Routes возвращает маршруты фичи users.
// Регистрация пользователей вынесена в auth feature (POST /auth/register).
func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{}
}

func encodeJSON(w http.ResponseWriter, v any) error {
	return json.NewEncoder(w).Encode(v)
}
