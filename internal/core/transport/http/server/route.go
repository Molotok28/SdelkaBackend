package core_http_server

import "net/http"

// Route описывает один HTTP-маршрут с опциональными middleware, применяемыми только к нему.
type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []func(http.Handler) http.Handler
}
