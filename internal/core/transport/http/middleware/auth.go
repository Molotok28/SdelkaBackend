package core_http_middleware

import (
	"net/http"

	core_auth "github.com/nikolaimoiseev01/SdelkaBackend/internal/core/auth"
)

// Auth проверяет JWT-токен из httpOnly-cookie и кладёт UserID в контекст.
// При ошибке возвращает 401 Unauthorized.
func Auth(jwtConfig core_auth.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(core_auth.TokenCookieName)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			claims, err := core_auth.ValidateToken(cookie.Value, jwtConfig)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := core_auth.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
