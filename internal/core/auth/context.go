package core_auth

import "context"

type userIDKey struct{}

// WithUserID добавляет ID пользователя в контекст после успешной аутентификации.
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFromContext извлекает ID пользователя из контекста.
// Второй аргумент false означает, что пользователь не аутентифицирован.
func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey{}).(int)
	return id, ok
}
