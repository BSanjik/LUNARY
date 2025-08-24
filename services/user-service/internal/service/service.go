// Бизнес-логика (профили, роли)
package service

import (
	"strings"
	"user-service/internal/model"
	"user-service/internal/storage"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/go-playground/validator/v10"
)

type UserService struct {
	Storage   *storage.Storage
	Validator *validator.Validate
}

func NewUserService(storage *storage.Storage) *UserService {
	return &UserService{
		Storage:   storage,
		Validator: validator.New(),
	}
}

// DTO для валидации создания профиля
type CreateProfileDTO struct {
	Name             *string `validate:"omitempty,max=50"`
	Lastname         *string `validate:"omitempty,max=50"`
	Sex              *string `validate:"omitempty,oneof=MALE FEMALE"`
	Profession       *string `validate:"omitempty,max=100"`
	City             *string `validate:"omitempty,max=100"`
	Language         *string `validate:"omitempty,oneof=RU KZ EN"`
	NotificationTone *string `validate:"omitempty,oneof=SOFT NEUTRAL STRICT"`
}

// DTO для валидации обновления профиля
type UpdateProfileDTO struct {
	Name             *string `validate:"omitempty,max=50"`
	Lastname         *string `validate:"omitempty,max=50"`
	Sex              *string `validate:"omitempty,oneof=MALE FEMALE"`
	Profession       *string `validate:"omitempty,max=100"`
	City             *string `validate:"omitempty,max=100"`
	Language         *string `validate:"omitempty,oneof=RU KZ EN"`
	NotificationTone *string `validate:"omitempty,oneof=SOFT NEUTRAL STRICT"`
}

// Валидация данных профиля
func (s *UserService) validateProfileData(req *model.UserProfileRequest, isUpdate bool) error {
	var dto interface{}

	if isUpdate {
		dto = &UpdateProfileDTO{
			Name:             req.Name,
			Lastname:         req.Lastname,
			Sex:              req.Sex,
			Profession:       req.Profession,
			City:             req.City,
			Language:         req.Language,
			NotificationTone: req.NotificationTone,
		}
	} else {
		dto = &CreateProfileDTO{
			Name:             req.Name,
			Lastname:         req.Lastname,
			Sex:              req.Sex,
			Profession:       req.Profession,
			City:             req.City,
			Language:         req.Language,
			NotificationTone: req.NotificationTone,
		}
	}

	if err := s.Validator.Struct(dto); err != nil {
		return s.formatValidationError(err)
	}

	return nil
}

func (s *UserService) formatValidationError(err error) error {
	var details []string
	for _, validErr := range err.(validator.ValidationErrors) {
		details = append(details, s.getValidationMessage(validErr))
	}
	return shared.NewValidationError("Ошибка валидации: " + strings.Join(details, "; "))
}

func (s *UserService) getValidationMessage(err validator.FieldError) string {
	fieldName := err.Field()
	switch err.Tag() {
	case "required":
		return fieldName + " обязательно для заполнения"
	case "max":
		return fieldName + " превышает максимальную длину"
	case "oneof":
		return fieldName + " имеет недопустимое значение"
	default:
		return fieldName + " имеет неверный формат"
	}
}

// Нормализация строковых полей
func (s *UserService) normalizeRequest(req *model.UserProfileRequest) {
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if req.Lastname != nil {
		trimmed := strings.TrimSpace(*req.Lastname)
		req.Lastname = &trimmed
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

// Создать профиль пользователя
func (s *UserService) CreateProfile(userID int64, req *model.UserProfileRequest) (*model.UserProfile, error) {
	// Нормализация данных
	s.normalizeRequest(req)

	// Валидация входных данных
	if err := s.validateProfileData(req, false); err != nil {
		shared.LogInfo("Profile creation validation failed", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}

	// Проверяем существование пользователя
	exists, err := s.Storage.UserExists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		shared.LogInfo("User not found for profile creation", map[string]interface{}{
			"user_id": userID,
		})
		return nil, shared.NewNotFoundError("Пользователь не найден")
	}

	// Создаем профиль
	profile, err := s.Storage.CreateProfile(userID, req)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Получить полный профиль пользователя
func (s *UserService) GetProfileByUserID(userID int64) (*model.FullUserProfile, error) {
	profile, err := s.Storage.GetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Обновить профиль пользователя
func (s *UserService) UpdateProfile(userID int64, req *model.UserProfileRequest) (*model.UserProfile, error) {
	// Нормализация данных
	s.normalizeRequest(req)

	// Валидация входных данных
	if err := s.validateProfileData(req, true); err != nil {
		shared.LogInfo("Profile update validation failed", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, err
	}

	// Обновляем профиль
	profile, err := s.Storage.UpdateProfile(userID, req)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// Удалить профиль пользователя
func (s *UserService) DeleteProfile(userID int64) error {
	if err := s.Storage.DeleteProfile(userID); err != nil {
		return err
	}

	return nil
}
