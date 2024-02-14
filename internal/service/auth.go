package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/sha3"
	"main.go/internal/entity"
	"main.go/internal/repository"
)

const (
	salt       = "57ct480 (T^&(6 n79TG789)"
	sighingKey = "v9R^V(&6&V*r^8^"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

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

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{time.Now().Add(tokenTTL)},
		IssuedAt:  &jwt.NumericDate{time.Now()},
	}, user.Id})

	return token.SignedString([]byte(sighingKey))
}

func generatePasswordHash(password string) string {
	hash := sha3.New256()
	hash.Write(([]byte(password)))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
