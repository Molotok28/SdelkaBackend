package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		fmt.Println("Failed to decode request body:", err)
	}

	if err := requestValidator.Struct(dest); err != nil {
		fmt.Println("Failed to validate request body:", err)
	}

	return nil
}
