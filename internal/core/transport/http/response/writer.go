package core_http_response

import "net/http"

var StatusCodeUninitialized = -1

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// WriteHeader перехватывает явный вызов WriteHeader от обработчика.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// Write перехватывает запись тела ответа и фиксирует неявный статус 200,
// который Go устанавливает автоматически при первом вызове Write без WriteHeader.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == StatusCodeUninitialized {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// GetStatusCode возвращает зафиксированный статус ответа или 200, если ответ ещё не был отправлен.
func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}
	return rw.statusCode
}
