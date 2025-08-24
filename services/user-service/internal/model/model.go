// Модели данных пользователя
package model

import "time"

type Sex string

const (
	SexMale   Sex = "MALE"
	SexFemale Sex = "FEMALE"
)

type NotificationTone string

const (
	ToneSoft    NotificationTone = "SOFT"
	ToneNeutral NotificationTone = "NEUTRAL"
	ToneStrict  NotificationTone = "STRICT"
)

// Модель профиля пользователя
type UserProfile struct {
	ID                      int64            `json:"id" gorm:"primaryKey"`
	UserID                  int64            `json:"user_id" gorm:"not null;uniqueIndex"`
	Name                    *string          `json:"name"`
	Lastname                *string          `json:"lastname"`
	Sex                     *Sex             `json:"sex"`
	Profession              *string          `json:"profession"`
	City                    string           `json:"city" gorm:"default:'Almaty'"`
	Timezone                string           `json:"timezone" gorm:"default:'Asia/Almaty'"`
	Language                string           `json:"language" gorm:"default:'RU'"`
	NotificationTone        NotificationTone `json:"notification_tone" gorm:"default:'NEUTRAL'"`
	LastSpecialScenarioDate *time.Time       `json:"last_special_scenario_date"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
}

// DTO для создания/обновления профиля
type UserProfileRequest struct {
	Name             *string `json:"name,omitempty"`
	Lastname         *string `json:"lastname,omitempty"`
	Sex              *string `json:"sex,omitempty"`
	Profession       *string `json:"profession,omitempty"`
	City             *string `json:"city,omitempty"`
	Language         *string `json:"language,omitempty"`
	NotificationTone *string `json:"notification_tone,omitempty"`
}

// Полный профиль для API ответов (объединение User + UserProfile)
type FullUserProfile struct {
	ID                      int64            `json:"id"`
	Phone                   string           `json:"phone"`                      // Из users
	Email                   *string          `json:"email"`                      // Из users
	Name                    *string          `json:"name"`                       // Из user_profiles
	Lastname                *string          `json:"lastname"`                   // Из user_profiles
	Sex                     *Sex             `json:"sex"`                        // Из user_profiles
	Profession              *string          `json:"profession"`                 // Из user_profiles
	City                    string           `json:"city"`                       // Из user_profiles
	Timezone                string           `json:"timezone"`                   // Из user_profiles
	Language                string           `json:"language"`                   // Из user_profiles
	NotificationTone        NotificationTone `json:"notification_tone"`          // Из user_profiles
	LastSpecialScenarioDate *time.Time       `json:"last_special_scenario_date"` // Из user_profiles
	CreatedAt               time.Time        `json:"created_at"`                 // Из users
}

// Ответ при работе с профилем
type ProfileResponse struct {
	Message string           `json:"message"`
	Profile *FullUserProfile `json:"profile"`
}
