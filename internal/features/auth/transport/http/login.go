package auth_transport_http

import (
	"errors"
	"net/http"
	"time"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_request "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/request"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	auth_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/auth/service"
	"go.uber.org/zap"
)

type loginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

func (h *AuthHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	var req loginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		log.Warn("invalid login request", zap.Error(err))
		core_http_response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Login(ctx, auth_service.LoginInput{
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	})
	if err != nil {
		if errors.Is(err, auth_service.ErrInvalidCredentials) {
			core_http_response.WriteError(w, http.StatusUnauthorized, "invalid phone number or password")
			return
		}
		log.Error("login failed", zap.Error(err))
		core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     core_auth.TokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.jwtConfig.TTL / time.Second),
	})

	if err := core_http_response.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		log.Error("write login response", zap.Error(err))
	}
}
