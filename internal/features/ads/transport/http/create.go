package ads_transport_http

import (
	"net/http"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_request "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/request"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
	"go.uber.org/zap"
)

type createAdRequest struct {
	Title       string  `json:"title" validate:"required,min=5,max=200"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Price       float64 `json:"price" validate:"gte=0"`
}

func (h *AdsHTTPHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	userID, ok := core_auth.UserIDFromContext(ctx)
	if !ok {
		core_http_response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createAdRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		log.Warn("invalid create ad request", zap.Error(err))
		core_http_response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ad, err := h.service.CreateAd(ctx, ads_service.CreateAdInput{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		UserID:      userID,
	})
	if err != nil {
		log.Error("create ad failed", zap.Error(err))
		core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := core_http_response.WriteJSON(w, http.StatusCreated, domainToResponse(ad)); err != nil {
		log.Error("write create ad response", zap.Error(err))
	}
}
