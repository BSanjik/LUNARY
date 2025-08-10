package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserContextKey = ctxKey("user")

func JWTMiddleware(secret []byte, exemptPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Освобождаемые маршруты
			if isExempt(r.URL.Path, exemptPaths) {
				shared.LogInfo("Запрос не требует авторизации", map[string]interface{}{
					"path":   r.URL.Path,
					"method": r.Method,
				})
				next.ServeHTTP(w, r)
				return
			}

			// Проверка заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				shared.LogWarn("Отсутствует токен авторизации", map[string]interface{}{
					"path":        r.URL.Path,
					"remote_addr": r.RemoteAddr,
					"user_agent":  r.UserAgent(),
				})
				shared.WriteErrorResponse(w, shared.NewAuthError("Требуется авторизация"))
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				shared.LogWarn("Неверный формат токена", map[string]interface{}{
					"path":               r.URL.Path,
					"remote_addr":        r.RemoteAddr,
					"auth_header_prefix": strings.Split(authHeader, " ")[0],
				})
				shared.WriteErrorResponse(w, shared.NewAuthError("Неверный формат токена авторизации"))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Парсинг токена
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					shared.LogWarn("Неожиданный алгоритм подписи JWT", map[string]interface{}{
						"algorithm": token.Header["alg"],
						"path":      r.URL.Path,
					})
					return nil, errors.New("неожиданный алгоритм подписи")
				}
				return secret, nil
			})

			if err != nil {
				shared.LogWarn("Ошибка парсинга JWT токена", map[string]interface{}{
					"error":       err.Error(),
					"path":        r.URL.Path,
					"remote_addr": r.RemoteAddr,
				})
				shared.WriteErrorResponse(w, shared.NewAuthError("Недействительный токен авторизации"))
				return
			}

			if !token.Valid {
				shared.LogWarn("Недействительный JWT токен", map[string]interface{}{
					"path":        r.URL.Path,
					"remote_addr": r.RemoteAddr,
				})
				shared.WriteErrorResponse(w, shared.NewAuthError("Токен авторизации истёк или недействителен"))
				return
			}

			// Извлечение claims из токена
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				shared.LogError(errors.New("invalid claims type"), "Не удалось извлечь claims из токена", map[string]interface{}{
					"path": r.URL.Path,
				})
				shared.WriteErrorResponse(w, shared.NewAuthError("Ошибка обработки токена"))
				return
			}

			// Проверка времени истечения (exp)
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					shared.LogWarn("Истёкший JWT токен", map[string]interface{}{
						"exp":          int64(exp),
						"current_time": time.Now().Unix(),
						"path":         r.URL.Path,
					})
					shared.WriteErrorResponse(w, shared.NewAuthError("Токен авторизации истёк"))
					return
				}
			}

			// Успешная авторизация - добавляем user_id в контекст
			ctx := r.Context()
			if sub, ok := claims["sub"].(string); ok {
				ctx = shared.WithUserID(ctx, sub)
			}

			// Добавляем claims в контекст
			ctx = context.WithValue(ctx, UserContextKey, claims)

			// Логируем с контекстом (request_id и user_id автоматически добавятся)
			shared.LogInfoWithContext(ctx, "Успешная авторизация", map[string]interface{}{
				"path":        r.URL.Path,
				"duration_ms": time.Since(start).Milliseconds(),
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Проверка, освобождён ли путь от проверки токена
func isExempt(path string, exemptPaths []string) bool {
	for _, prefix := range exemptPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
