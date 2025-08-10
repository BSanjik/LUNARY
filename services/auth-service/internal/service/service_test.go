package service

import (
	"auth-service/internal/storage"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *storage.Storage) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	storage := &storage.Storage{DB: db}
	return db, mock, storage
}

func TestAuthService_Register_Success(t *testing.T) {
	db, mock, storage := setupMockDB(t)
	defer db.Close()

	service := &AuthService{Storage: storage}

	mock.ExpectQuery("SELECT EXISTS").WithArgs("+77085888664", sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "email", "created_at"}).AddRow(1, "+77085888664", nil, "2024-01-01"))

	response, err := service.Register("+77085888664", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "Пользователь успешно зарегистрирован", response.Message)
	assert.Equal(t, "+77085888664", response.User.Phone)

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestAuthService_Register_UserExists(t *testing.T) {
	db, mock, storage := setupMockDB(t)
	defer db.Close()

	service := &AuthService{Storage: storage}

	// Mock проверка существования пользователя - пользователь существует
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("+77085888664").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Выполняем тест
	response, err := service.Register("+77085888664", "password123")

	// Проверяем результат
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "пользователь с таким номером уже существует")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthService_Login_Success(t *testing.T) {
	db, mock, storage := setupMockDB(t)
	defer db.Close()

	service := &AuthService{Storage: storage}

	// Хешированный пароль для "password123"
	hashedPassword := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"

	// Mock поиск пользователя
	mock.ExpectQuery("SELECT id, phone, email, password, created_at FROM users").
		WithArgs("+77085888664").
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone", "email", "password", "created_at"}).
			AddRow(1, "+77085888664", nil, hashedPassword, "2024-01-01"))

	// Выполняем тест
	response, err := service.Login("+77085888664", "password123")

	// Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "Авторизация успешна", response.Message)

	assert.NoError(t, mock.ExpectationsWereMet())
}
