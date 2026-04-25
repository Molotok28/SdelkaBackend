package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

// DecodeAndValidateRequest декодирует JSON-тело запроса в dest и валидирует его по тегам validate.
// Возвращает ошибку, если декодирование или валидация не прошли.
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode request body: %w", err)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf("validate request body: %w", err)
	}

	return nil
}
