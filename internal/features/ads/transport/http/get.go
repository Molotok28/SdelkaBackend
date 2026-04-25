package ads_transport_http

import (
	"errors"
	"net/http"
	"strconv"

	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
	"go.uber.org/zap"
)

func (h *AdsHTTPHandler) GetAd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		core_http_response.WriteError(w, http.StatusBadRequest, "invalid ad id")
		return
	}

	ad, err := h.service.GetAd(ctx, id)
	if err != nil {
		if errors.Is(err, ads_service.ErrNotFound) {
			core_http_response.WriteError(w, http.StatusNotFound, "ad not found")
			return
		}
		log.Error("get ad failed", zap.Int("id", id), zap.Error(err))
		core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := core_http_response.WriteJSON(w, http.StatusOK, domainToResponse(ad)); err != nil {
		log.Error("write get ad response", zap.Error(err))
	}
}
