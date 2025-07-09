package infrastructure

import (
	"context"
	"fmt"
	"strings"

	uuid "github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(username string) (string, error)
	ValidateToken(tokenString string) (string, error)
}

type tokenService struct {
	redisService RedisService
}

func NewTokenService(redisService RedisService) TokenService {
	return &tokenService{redisService: redisService}
}

// generates a new token
func (s *tokenService) GenerateToken(username string) (string, error) {
	token := uuid.New().String() + ":" + username
	return token, nil
}

// validates a token
func (s *tokenService) ValidateToken(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format")
	}

	username := parts[1]
	if username == "" {
		return "", fmt.Errorf("invalid token format")
	}
	
	uuidPart := parts[0]
	if _, err := uuid.Parse(uuidPart); err != nil {
		return "", fmt.Errorf("invalid token format")
	}

	// Check if the token exists in Redis
	exists, err := s.redisService.GetClient().Exists(context.TODO(), "session:"+tokenString).Result()
	if err != nil {
		return "", fmt.Errorf("failed to check token existence: %v", err)
	}

	if exists == 0 {
		return "", fmt.Errorf("token does not exist")
	}

	return username, nil
}