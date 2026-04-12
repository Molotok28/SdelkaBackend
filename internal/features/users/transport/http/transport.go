package users_transport_http

import (
	"net/http"

	core_http_server "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	userService UsersService
}

type UsersService interface {
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
