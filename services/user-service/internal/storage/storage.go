// Доступ к БД (PostgreSQL)
package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"user-service/internal/model"

	"github.com/BSanjik/LUNARY/services/common/shared"
	_ "github.com/lib/pq"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	AuthServiceURL string
}

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg *Config) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		shared.LogError(err, "Failed to open database connection", map[string]interface{}{
			"host": cfg.DBHost,
			"port": cfg.DBPort,
			"db":   cfg.DBName,
		})
		return nil, err
	}

	// Настройка connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Connection check with retries
	for i := 1; i <= 10; i++ {
		err := db.Ping()
		if err == nil {
			shared.LogInfo("Successfully connected to database", map[string]interface{}{
				"host": cfg.DBHost,
				"port": cfg.DBPort,
				"db":   cfg.DBName,
			})
			return &Storage{DB: db}, nil
		}
		shared.LogError(err, fmt.Sprintf("Database ping attempt %d failed", i), map[string]interface{}{
			"attempt": i,
			"host":    cfg.DBHost,
		})
		time.Sleep(2 * time.Second)
	}

	return &Storage{DB: db}, nil
}

// Создать профиль пользователя
func (s *Storage) CreateProfile(userID int64, req *model.UserProfileRequest) (*model.UserProfile, error) {
	query := `
		INSERT INTO user_profiles (
			user_id, name, lastname, sex, profession, city, 
			timezone, language, notification_tone
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, name, lastname, sex, profession, 
				  city, timezone, language, notification_tone, 
				  last_special_scenario_date, created_at, updated_at
	`

	profile := &model.UserProfile{}

	err := s.DB.QueryRow(
		query,
		userID,
		req.Name,
		req.Lastname,
		req.Sex,
		req.Profession,
		getStringValue(req.City, "Almaty"),
		"Asia/Almaty", // default timezone
		getStringValue(req.Language, "RU"),
		getStringValue(req.NotificationTone, "NEUTRAL"),
	).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Name,
		&profile.Lastname,
		&profile.Sex,
		&profile.Profession,
		&profile.City,
		&profile.Timezone,
		&profile.Language,
		&profile.NotificationTone,
		&profile.LastSpecialScenarioDate,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		shared.LogError(err, "Failed to create user profile", map[string]interface{}{
			"user_id": userID,
		})
		return nil, shared.NewInternalError("Ошибка создания профиля")
	}

	shared.LogInfo("User profile created successfully", map[string]interface{}{
		"user_id":    userID,
		"profile_id": profile.ID,
	})

	return profile, nil
}

// Получить полный профиль по user_id (объединение users + user_profiles)
func (s *Storage) GetProfileByUserID(userID int64) (*model.FullUserProfile, error) {
	query := `
		SELECT 
			u.id, u.phone, u.email, u.created_at,
			up.name, up.lastname, up.sex, up.profession, 
			up.city, up.timezone, up.language, up.notification_tone,
			up.last_special_scenario_date
		FROM users u
		LEFT JOIN user_profiles up ON u.id = up.user_id
		WHERE u.id = $1
	`

	profile := &model.FullUserProfile{}

	err := s.DB.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.Phone,
		&profile.Email,
		&profile.CreatedAt,
		&profile.Name,
		&profile.Lastname,
		&profile.Sex,
		&profile.Profession,
		&profile.City,
		&profile.Timezone,
		&profile.Language,
		&profile.NotificationTone,
		&profile.LastSpecialScenarioDate,
	)

	if err == sql.ErrNoRows {
		shared.LogInfo("User profile not found", map[string]interface{}{
			"user_id": userID,
		})
		return nil, shared.NewNotFoundError("Профиль пользователя не найден")
	} else if err != nil {
		shared.LogError(err, "Failed to get user profile", map[string]interface{}{
			"user_id": userID,
		})
		return nil, shared.NewInternalError("Ошибка получения профиля")
	}

	return profile, nil
}

// Обновить профиль
func (s *Storage) UpdateProfile(userID int64, req *model.UserProfileRequest) (*model.UserProfile, error) {
	// Динамически строим UPDATE запрос только для переданных полей
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}
	if req.Lastname != nil {
		setParts = append(setParts, fmt.Sprintf("lastname = $%d", argIndex))
		args = append(args, req.Lastname)
		argIndex++
	}
	if req.Sex != nil {
		setParts = append(setParts, fmt.Sprintf("sex = $%d", argIndex))
		args = append(args, *req.Sex)
		argIndex++
	}
	if req.Profession != nil {
		setParts = append(setParts, fmt.Sprintf("profession = $%d", argIndex))
		args = append(args, req.Profession)
		argIndex++
	}
	if req.City != nil {
		setParts = append(setParts, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *req.City)
		argIndex++
	}
	if req.Language != nil {
		setParts = append(setParts, fmt.Sprintf("language = $%d", argIndex))
		args = append(args, *req.Language)
		argIndex++
	}
	if req.NotificationTone != nil {
		setParts = append(setParts, fmt.Sprintf("notification_tone = $%d", argIndex))
		args = append(args, *req.NotificationTone)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, shared.NewValidationError("Нет полей для обновления")
	}

	// Добавляем updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Добавляем user_id для WHERE
	args = append(args, userID)

	setClause := strings.Join(setParts, ", ")

	query := fmt.Sprintf(`
		UPDATE user_profiles 
		SET %s
		WHERE user_id = $%d
		RETURNING id, user_id, name, lastname, sex, profession, 
				  city, timezone, language, notification_tone, 
				  last_special_scenario_date, created_at, updated_at
	`, setClause, argIndex)

	profile := &model.UserProfile{}

	err := s.DB.QueryRow(query, args...).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Name,
		&profile.Lastname,
		&profile.Sex,
		&profile.Profession,
		&profile.City,
		&profile.Timezone,
		&profile.Language,
		&profile.NotificationTone,
		&profile.LastSpecialScenarioDate,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		shared.LogInfo("Profile not found for update", map[string]interface{}{
			"user_id": userID,
		})
		return nil, shared.NewNotFoundError("Профиль не найден")
	} else if err != nil {
		shared.LogError(err, "Failed to update profile", map[string]interface{}{
			"user_id": userID,
		})
		return nil, shared.NewInternalError("Ошибка обновления профиля")
	}

	shared.LogInfo("Profile updated successfully", map[string]interface{}{
		"user_id":    userID,
		"profile_id": profile.ID,
	})

	return profile, nil
}

// Удалить профиль
func (s *Storage) DeleteProfile(userID int64) error {
	query := "DELETE FROM user_profiles WHERE user_id = $1"

	result, err := s.DB.Exec(query, userID)
	if err != nil {
		shared.LogError(err, "Failed to delete profile", map[string]interface{}{
			"user_id": userID,
		})
		return shared.NewInternalError("Ошибка удаления профиля")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		shared.LogError(err, "Failed to get rows affected", map[string]interface{}{
			"user_id": userID,
		})
		return shared.NewInternalError("Ошибка удаления профиля")
	}

	if rowsAffected == 0 {
		shared.LogInfo("Profile not found for deletion", map[string]interface{}{
			"user_id": userID,
		})
		return shared.NewNotFoundError("Профиль не найден")
	}

	shared.LogInfo("Profile deleted successfully", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// Проверить существование пользователя в таблице users
func (s *Storage) UserExists(userID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"

	err := s.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		shared.LogError(err, "Failed to check user existence", map[string]interface{}{
			"user_id": userID,
		})
		return false, shared.NewInternalError("Ошибка проверки пользователя")
	}

	return exists, nil
}

// Вспомогательная функция для получения значения строки или дефолта
func getStringValue(ptr *string, defaultValue string) string {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
