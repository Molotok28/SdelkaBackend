package auth_transport_http

import (
	"net/http"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

func (h *AuthHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	http.SetCookie(w, &http.Cookie{
		Name:     core_auth.TokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	if err := core_http_response.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		log.Error("write logout response", zap.Error(err))
	}
}
