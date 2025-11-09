package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Romasmi/social-network/internal/services"
	"github.com/Romasmi/social-network/internal/utils"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	SessionService *services.SessionService
}

func CreateAuthMiddleware(sessionService *services.SessionService) *AuthMiddleware {
	return &AuthMiddleware{
		SessionService: sessionService,
	}
}

func (m *AuthMiddleware) Process(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get bearer token
		authHeader := r.Header.Get("Authorization")
		authFields := strings.Fields(authHeader)
		if len(authFields) != 2 || authFields[0] != "Bearer" {
			utils.JsonError(w, http.StatusForbidden, nil)
			return
		}
		token := authFields[1]
		sessionId, err := uuid.Parse(token)
		if err != nil {
			utils.ErrorInvalidRequestBody(w, fmt.Errorf("invalid token"))
			return
		}

		if isValidSession, err := m.SessionService.IsValidSession(r.Context(), sessionId); isValidSession == false || err != nil {
			if err != nil {
				fmt.Printf("error while check session: %v\n", err)
			}
			utils.JsonError(w, http.StatusForbidden, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
