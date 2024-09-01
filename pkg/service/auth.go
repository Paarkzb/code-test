package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"medodstest/internal/model"
	"medodstest/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	salt       = "kjasdhflkqwurh"
	signingKey = "askdjfsa;ldfkjdsal;128"
	tokenTTL   = 24 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"sub"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GetUser(userId int) (model.UserResponse, error) {
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&tokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.Id,
		})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Неправильный метод подписи")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, nil
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("Тип токена не совпадает с типом tokenClaims")
	}

	return claims.UserId, nil
}
