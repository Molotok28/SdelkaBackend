package ads_transport_http

import (
	"errors"
	"net/http"
	"strconv"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_request "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/request"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
	"go.uber.org/zap"
)

type updateAdRequest struct {
	Title       string  `json:"title" validate:"required,min=5,max=200"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Price       float64 `json:"price" validate:"gte=0"`
}

func (h *AdsHTTPHandler) UpdateAd(w http.ResponseWriter, r *http.Request) {
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

	var req updateAdRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		log.Warn("invalid update ad request", zap.Error(err))
		core_http_response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ad, err := h.service.UpdateAd(ctx, ads_service.UpdateAdInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		RequesterID: userID,
	})
	if err != nil {
		switch {
		case errors.Is(err, ads_service.ErrNotFound):
			core_http_response.WriteError(w, http.StatusNotFound, "ad not found")
		case errors.Is(err, ads_service.ErrForbidden):
			core_http_response.WriteError(w, http.StatusForbidden, "you don't own this ad")
		default:
			log.Error("update ad failed", zap.Int("id", id), zap.Error(err))
			core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	if err := core_http_response.WriteJSON(w, http.StatusOK, domainToResponse(ad)); err != nil {
		log.Error("write update ad response", zap.Error(err))
	}
}
