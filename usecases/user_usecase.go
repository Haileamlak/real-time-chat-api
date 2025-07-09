package usecases

import (
	"context"
	"errors"

	"github.com/haileamlak/chat-system/infrastructure"
	"github.com/haileamlak/chat-system/repositories"
)

type UserUseCase interface {
	Register(ctx context.Context, username string, password string) error
	Login(ctx context.Context, username string, password string) (string, error)
}

type userUseCase struct {
	userRepo      repositories.UserRepository
	passwordService infrastructure.PasswordService
	tokenService    infrastructure.TokenService
}

func NewUserUseCase(userRepo repositories.UserRepository, passwordService infrastructure.PasswordService, tokenService infrastructure.TokenService) UserUseCase {
	return &userUseCase{
		userRepo:      userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

func (u *userUseCase) Register(ctx context.Context, username string, password string) error {
	// Check if user already exists
	exists, err := u.userRepo.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := u.passwordService.HashPassword(password)
	if err != nil {
		return err
	}

	// Create the user
	return u.userRepo.CreateUser(ctx, username, hashedPassword)
}

func (u *userUseCase) Login(ctx context.Context, username string, password string) (string, error) {
	// Get the user's hashed password
	hashedPassword, err := u.userRepo.GetUserPassword(ctx, username)
	if err != nil {
		return "", err
	}

	// Compare the passwords
	if err := u.passwordService.ComparePasswords(hashedPassword, password); err != nil {
		return "", err
	}

	// Generate a new session token
	token, err := u.tokenService.GenerateToken(username)
	if err != nil {
		return "", err
	}

	// Store the session token in the repository
	if err := u.userRepo.SaveSession(ctx, token, username); err != nil {
		return "", err
	}

	return token, nil
}	