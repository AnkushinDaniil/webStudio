package service

import (
	"fmt"

	"golang.org/x/crypto/sha3"
	"main.go/internal/entity"
	"main.go/internal/repository"
)

const SALT = "57ct480 (T^&(6 n79TG789)"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha3.New256()
	hash.Write(([]byte(password)))

	return fmt.Sprintf("%x", hash.Sum([]byte(SALT)))
}
