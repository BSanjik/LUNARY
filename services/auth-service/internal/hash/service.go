package hash

import "golang.org/x/crypto/bcrypt"

type HashService struct{}

func New() *HashService {
	return &HashService{}
}

func (h *HashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (h *HashService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
