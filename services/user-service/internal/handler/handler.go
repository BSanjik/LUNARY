// HTTP обработчики профиля пользователя
package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"user-service/internal/model"
	"user-service/internal/service"
	"user-service/internal/token"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

type Handler struct {
	Service *service.UserService
}

type CreateProfileRequest struct {
	Name             *string `json:"name,omitempty"`
	Lastname         *string `json:"lastname,omitempty"`
	Sex              *string `json:"sex,omitempty"`
	Profession       *string `json:"profession,omitempty"`
	City             *string `json:"city,omitempty"`
	Language         *string `json:"language,omitempty"`
	NotificationTone *string `json:"notification_tone,omitempty"`
}

type UpdateProfileRequest struct {
	Name             *string `json:"name,omitempty"`
	Lastname         *string `json:"lastname,omitempty"`
	Sex              *string `json:"sex,omitempty"`
	Profession       *string `json:"profession,omitempty"`
	City             *string `json:"city,omitempty"`
	Language         *string `json:"language,omitempty"`
	NotificationTone *string `json:"notification_tone,omitempty"`
}

// ErrorResponse структура для ошибок
type ErrorResponse struct {
	Error shared.AppError `json:"error"`
}

const maxRequestSize = 1 << 20 // 1MB

// Простая функция для декодирования JSON
func (h *Handler) decodeRequest(r *http.Request, dst interface{}) error {
	// Лимит размера запроса
	r.Body = http.MaxBytesReader(nil, r.Body, maxRequestSize)

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return shared.NewValidationError("Неверный JSON")
	}

	// Применяем trim для строковых полей
	h.trimStringFields(dst)

	return nil
}

func (h *Handler) trimStringFields(v interface{}) {
	switch req := v.(type) {
	case *CreateProfileRequest:
		if req.Name != nil {
			trimmed := strings.TrimSpace(*req.Name)
			req.Name = &trimmed
		}
		if req.Lastname != nil {
			trimmed := strings.TrimSpace(*req.Lastname)
			req.Lastname = &trimmed
		}
		if req.Sex != nil {
			trimmed := strings.TrimSpace(*req.Sex)
			req.Sex = &trimmed
		}
		if req.Profession != nil {
			trimmed := strings.TrimSpace(*req.Profession)
			req.Profession = &trimmed
		}
		if req.City != nil {
			trimmed := strings.TrimSpace(*req.City)
			req.City = &trimmed
		}
		if req.Language != nil {
			trimmed := strings.TrimSpace(*req.Language)
			req.Language = &trimmed
		}
		if req.NotificationTone != nil {
			trimmed := strings.TrimSpace(*req.NotificationTone)
			req.NotificationTone = &trimmed
		}
	case *UpdateProfileRequest:
		if req.Name != nil {
			trimmed := strings.TrimSpace(*req.Name)
			req.Name = &trimmed
		}
		if req.Lastname != nil {
			trimmed := strings.TrimSpace(*req.Lastname)
			req.Lastname = &trimmed
		}
		if req.Sex != nil {
			trimmed := strings.TrimSpace(*req.Sex)
			req.Sex = &trimmed
		}
		if req.Profession != nil {
			trimmed := strings.TrimSpace(*req.Profession)
			req.Profession = &trimmed
		}
		if req.City != nil {
			trimmed := strings.TrimSpace(*req.City)
			req.City = &trimmed
		}
		if req.Language != nil {
			trimmed := strings.TrimSpace(*req.Language)
			req.Language = &trimmed
		}
		if req.NotificationTone != nil {
			trimmed := strings.TrimSpace(*req.NotificationTone)
			req.NotificationTone = &trimmed
		}
	}
}

// Преобразование DTO в модель
func (h *Handler) createRequestToModel(req *CreateProfileRequest) *model.UserProfileRequest {
	return &model.UserProfileRequest{
		Name:             req.Name,
		Lastname:         req.Lastname,
		Sex:              req.Sex,
		Profession:       req.Profession,
		City:             req.City,
		Language:         req.Language,
		NotificationTone: req.NotificationTone,
	}
}

func (h *Handler) updateRequestToModel(req *UpdateProfileRequest) *model.UserProfileRequest {
	return &model.UserProfileRequest{
		Name:             req.Name,
		Lastname:         req.Lastname,
		Sex:              req.Sex,
		Profession:       req.Profession,
		City:             req.City,
		Language:         req.Language,
		NotificationTone: req.NotificationTone,
	}
}

// Middleware для авторизации
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			shared.WriteErrorResponse(w, shared.NewAuthError("Authorization header is required"))
			return
		}

		// Проверяем формат "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			shared.WriteErrorResponse(w, shared.NewAuthError("Authorization header must be in format: Bearer <token>"))
			return
		}

		tokenString := tokenParts[1]

		// Парсим токен
		claims, err := token.ParseToken(tokenString)
		if err != nil {
			shared.WriteErrorResponse(w, shared.NewAuthError("Invalid token: "+err.Error()))
			return
		}

		// Сохраняем данные пользователя в контексте
		r = r.WithContext(token.ContextWithUserID(r.Context(), claims.UserID))
		r = r.WithContext(token.ContextWithPhone(r.Context(), claims.Phone))

		next.ServeHTTP(w, r)
	}
}

// Создать профиль пользователя
func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := token.UserIDFromContext(r.Context())
	if err != nil {
		shared.WriteErrorResponse(w, shared.NewAuthError("User ID not found in context"))
		return
	}

	var req CreateProfileRequest
	if err := h.decodeRequest(r, &req); err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	modelReq := h.createRequestToModel(&req)
	_, err = h.Service.CreateProfile(userID, modelReq)
	if err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	response := model.ProfileResponse{
		Message: "Профиль создан успешно",
		Profile: nil, // Можно вернуть полный профиль если нужно
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Получить профиль пользователя
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := token.UserIDFromContext(r.Context())
	if err != nil {
		shared.WriteErrorResponse(w, shared.NewAuthError("User ID not found in context"))
		return
	}

	profile, err := h.Service.GetProfileByUserID(userID)
	if err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	response := model.ProfileResponse{
		Message: "Профиль получен успешно",
		Profile: profile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Обновить профиль пользователя
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := token.UserIDFromContext(r.Context())
	if err != nil {
		shared.WriteErrorResponse(w, shared.NewAuthError("User ID not found in context"))
		return
	}

	var req UpdateProfileRequest
	if err := h.decodeRequest(r, &req); err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	modelReq := h.updateRequestToModel(&req)
	_, err = h.Service.UpdateProfile(userID, modelReq)
	if err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	response := model.ProfileResponse{
		Message: "Профиль обновлен успешно",
		Profile: nil, // Можно вернуть обновленный полный профиль
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Удалить профиль пользователя
func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := token.UserIDFromContext(r.Context())
	if err != nil {
		shared.WriteErrorResponse(w, shared.NewAuthError("User ID not found in context"))
		return
	}

	if err := h.Service.DeleteProfile(userID); err != nil {
		shared.WriteErrorResponse(w, err)
		return
	}

	response := map[string]string{
		"message": "Профиль удален успешно",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
