package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/google/uuid"
)

type SessionService struct {
	sessionRepo *repository.SessionRepository
}

func CreateSessionService(sessionRepo *repository.SessionRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) CreateSession(ctx context.Context, userId uuid.UUID, metadata json.RawMessage, ttl time.Duration) (*models.Session, error) {
	if ttl <= 0 {
		return nil, errors.New("ttl must be greater than 0")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	session := &models.Session{
		ID:        id,
		UserId:    userId,
		Metadata:  metadata,
		ExpiresAt: time.Now().Add(ttl),
	}
	return s.sessionRepo.CreateSession(ctx, session)
}

func (s *SessionService) GetSessionById(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	return s.sessionRepo.GetSessionById(ctx, id)
}
