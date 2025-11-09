package services

import (
	"context"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	InvalidCredentialsError = errors.New("invalid credentials")
)

type AuthService struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

func CreateAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, profileId uuid.UUID, passwordHash string) (*models.Session, error) {
	user, err := s.UserRepo.GetUserByProfileId(ctx, profileId)
	if err != nil {
		// TODO process other types of errors
		return nil, InvalidCredentialsError
	}
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordHash))
	if err != nil {
		return nil, InvalidCredentialsError
	}

	sessionService := CreateSessionService(s.SessionRepo)
	session, err := sessionService.CreateSession(ctx, user.ID, nil, DefaultSessionTTL)
	if err != nil {
		return nil, err
	}
	return session, nil
}
