package services

import (
	"context"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func CreateAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, profileId uuid.UUID, passwordHash string) (*models.Session, error) {
	user, err := s.UserRepo.GetUserByProfileId(ctx, profileId)
	if err != nil {
		return nil, err
	}
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordHash))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
