package users_transport_http

import (
	"net/http"

	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_request "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/request"
	users_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/users/service"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Surname     string  `json:"surname" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15"`
}

type CreateUserResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	log.Debug("invoke CreateUser handler called")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		log.Warn("invalid CreateUser request", zap.Error(err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(ctx, users_service.CreateUserInput{
		Name:        request.Name,
		Surname:     request.Surname,
		PhoneNumber: request.PhoneNumber,
	})
	if err != nil {
		log.Error("failed to create user", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	response := CreateUserResponse{
		ID:          user.ID,
		Version:     user.Version,
		Name:        user.Name,
		Surname:     user.Surname,
		PhoneNumber: user.PhoneNumber,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	if err := encodeJSON(rw, response); err != nil {
		log.Error("failed to encode CreateUser response", zap.Error(err))
	}
}
