package ads_transport_http

import (
	"net/http"
	"strconv"

	core_logger "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/logger"
	core_http_response "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/transport/http/response"
	ads_service "github.com/nikolaimoiseev01/SdelkaBackend/internal/features/ads/service"
	"go.uber.org/zap"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type listAdsResponse struct {
	Items  []adResponse `json:"items"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

func (h *AdsHTTPHandler) ListAds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	limit := parseQueryInt(r, "limit", defaultLimit)
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}

	offset := parseQueryInt(r, "offset", 0)
	if offset < 0 {
		offset = 0
	}

	result, err := h.service.ListAds(ctx, ads_service.ListAdsInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Error("list ads failed", zap.Error(err))
		core_http_response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]adResponse, 0, len(result.Items))
	for _, ad := range result.Items {
		items = append(items, domainToResponse(ad))
	}

	resp := listAdsResponse{
		Items:  items,
		Total:  result.Total,
		Limit:  result.Limit,
		Offset: result.Offset,
	}

	if err := core_http_response.WriteJSON(w, http.StatusOK, resp); err != nil {
		log.Error("write list ads response", zap.Error(err))
	}
}

func parseQueryInt(r *http.Request, key string, defaultVal int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}
	return v
}
