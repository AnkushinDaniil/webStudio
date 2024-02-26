package service

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/sha3"
	"main.go/internal/entity"
)

type AuthorizationRepository interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

const (
	salt       = "57ct480 (T^&(6 n79TG789)"
	sighingKey = "v9R^V(&6&V*r^8^"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type AuthorizationService struct {
	repo AuthorizationRepository
}

func NewAuthorizationService(repo AuthorizationRepository) *AuthorizationService {
	return &AuthorizationService{repo: repo}
}

func (s *AuthorizationService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthorizationService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := func() *jwt.Token {
		var (
			method jwt.SigningMethod = jwt.SigningMethodHS256
			claims jwt.Claims        = &tokenClaims{
				jwt.RegisteredClaims{
					Issuer:    "",
					Subject:   "",
					Audience:  []string{},
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(tokenTTL)},
					NotBefore: &jwt.NumericDate{
						Time: time.Time{},
					},
					IssuedAt: &jwt.NumericDate{Time: time.Now()},
					ID:       "",
				},
				user.ID,
			}
		)
		return &jwt.Token{
			Raw:       "",
			Method:    method,
			Header:    map[string]interface{}{"typ": "JWT", "alg": method.Alg()},
			Claims:    claims,
			Signature: []byte{},
			Valid:     false,
		}
	}()

	return token.SignedString([]byte(sighingKey))
}

func (s *AuthorizationService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:   "",
				Subject:  "",
				Audience: []string{},
				ExpiresAt: &jwt.NumericDate{
					Time: time.Time{},
				},
				NotBefore: &jwt.NumericDate{
					Time: time.Time{},
				},
				IssuedAt: &jwt.NumericDate{
					Time: time.Time{},
				},
				ID: "",
			},
			UserID: 0,
		},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			return []byte(sighingKey), nil
		})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not type *tokenClaims")
	}

	return claims.UserID, nil
}

func generatePasswordHash(password string) string {
	hash := sha3.New256()
	hash.Write(([]byte(password)))

	return hex.EncodeToString(hash.Sum([]byte(salt)))
}
