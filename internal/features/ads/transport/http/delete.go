package ads_transport_http

import (
	"errors"
	"net/http"
	"strconv"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
	"go.uber.org/zap"
)

func (h *AdsHTTPHandler) DeleteAd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		core_http_response.WriteError(w, http.StatusBadRequest, "invalid ad id")
		return
	}

	userID, ok := core_auth.UserIDFromContext(ctx)
	if !ok {
		core_http_response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.service.DeleteAd(ctx, id, userID); err != nil {
		switch {
		case errors.Is(err, ads_service.ErrNotFound):
			core_http_response.WriteError(w, http.StatusNotFound, "ad not found")
		case errors.Is(err, ads_service.ErrForbidden):
			core_http_response.WriteError(w, http.StatusForbidden, "you don't own this ad")
		default:
			log.Error("delete ad failed", zap.Int("id", id), zap.Error(err))
			core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	if err := core_http_response.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		log.Error("write delete ad response", zap.Error(err))
	}
}
