package auth_transport_http

import (
	"net/http"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

func (h *AuthHTTPHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	userID, ok := core_auth.UserIDFromContext(ctx)
	if !ok {
		core_http_response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		log.Error("get current user failed", zap.Error(err))
		core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := core_http_response.WriteJSON(w, http.StatusOK, userResponse{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		PhoneNumber: user.PhoneNumber,
	}); err != nil {
		log.Error("write me response", zap.Error(err))
	}
}
