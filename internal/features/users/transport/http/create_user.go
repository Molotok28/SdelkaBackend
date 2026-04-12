package users_transport_http

import (
	"encoding/json"
	"net/http"

	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
)

type CreateUserRequest struct {
	Name             string  `json:"name" validate:"required,min=3,max=100"`
	Surname          string  `json:"surname" validate:"required,min=3,max=100"`
	PhonePhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=100"`
}

type CreateUserResponse struct {
	ID               int     `json:"id"`
	Version          int     `json:"version"`
	Name             string  `json:"name"`
	Surname          string  `json:"surname"`
	PhonePhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	log.Debug("invoke CreateUser handler called")
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)

}
