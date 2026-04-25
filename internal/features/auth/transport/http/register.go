package auth_transport_http

import (
	"errors"
	"net/http"

	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_request "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/request"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	auth_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/service"
	"go.uber.org/zap"
)

type registerRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Surname     string `json:"surname" validate:"required,min=3,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

type userResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *AuthHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	var req registerRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		log.Warn("invalid register request", zap.Error(err))
		core_http_response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(ctx, auth_service.RegisterInput{
		Name:        req.Name,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		if errors.Is(err, auth_service.ErrPhoneAlreadyExists) {
			core_http_response.WriteError(w, http.StatusConflict, "phone number already registered")
			return
		}
		log.Error("register failed", zap.Error(err))
		core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := core_http_response.WriteJSON(w, http.StatusCreated, userResponse{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		PhoneNumber: user.PhoneNumber,
	}); err != nil {
		log.Error("write register response", zap.Error(err))
	}
}
