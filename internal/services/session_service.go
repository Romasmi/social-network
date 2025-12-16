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

const DefaultSessionTTL = 24 * time.Hour

type SessionService struct {
	sessionRepo *repository.SessionRepository
}

func CreateSessionService(sessionRepo *repository.SessionRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) CreateSession(ctx context.Context, userId uuid.UUID, metadata json.RawMessage, ttl time.Duration) (*models.Session, error) {
	// TODO invalidate previous sessions on new login
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

func (s *SessionService) IsValidSession(ctx context.Context, sessionId uuid.UUID) (bool, *models.Session, error) {
	session, err := s.sessionRepo.GetSessionById(ctx, sessionId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	if time.Now().After(session.ExpiresAt) {
		return false, nil, nil
	}
	return true, session, nil
}
