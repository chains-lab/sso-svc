package db

import (
	"encoding/json"
	"net/http"

	"github.com/cifra-city/comtools/httpkit"
	"github.com/cifra-city/sso-oauth/internal/data/db/sqlcore"
	"github.com/google/uuid"
)

type Sessions interface {
	Create(r *http.Request, userID uuid.UUID, token string, deviceName string, deviceData json.RawMessage) (sqlcore.Session, error)

	GetByID(r *http.Request, id uuid.UUID) (sqlcore.Session, error)
	GetSession(r *http.Request, id uuid.UUID, userID uuid.UUID) (sqlcore.Session, error)
	GetSessions(r *http.Request, userID uuid.UUID) ([]sqlcore.Session, error)

	GetToken(r *http.Request, id uuid.UUID, userID uuid.UUID) (string, error)
	UpdateToken(r *http.Request, id uuid.UUID, token string) error

	DeleteAll(r *http.Request, id uuid.UUID) error
	Delete(r *http.Request, id uuid.UUID, userID uuid.UUID) error
}

type sessions struct {
	queries *sqlcore.Queries
}

func NewSession(queries *sqlcore.Queries) Sessions {
	return &sessions{queries: queries}
}

func (s *sessions) Create(r *http.Request, userID uuid.UUID, token string, deviceName string, deviceData json.RawMessage) (sqlcore.Session, error) {
	client := httpkit.GetUserAgent(r)
	IP := httpkit.GetClientIP(r)
	return s.queries.CreateSession(r.Context(), sqlcore.CreateSessionParams{
		UserID:     userID,
		Token:      token,
		DeviceName: deviceName,
		Client:     client,
		Ip:         IP,
	})
}

func (s *sessions) GetByID(r *http.Request, id uuid.UUID) (sqlcore.Session, error) {
	return s.queries.GetSession(r.Context(), id)
}

func (s *sessions) GetSession(r *http.Request, id uuid.UUID, userID uuid.UUID) (sqlcore.Session, error) {
	return s.queries.GetUserSession(r.Context(), sqlcore.GetUserSessionParams{
		ID:     id,
		UserID: userID,
	})
}

func (s *sessions) GetSessions(r *http.Request, userID uuid.UUID) ([]sqlcore.Session, error) {
	return s.queries.GetSessionsByUserID(r.Context(), userID)
}

func (s *sessions) GetToken(r *http.Request, id uuid.UUID, userID uuid.UUID) (string, error) {
	return s.queries.GetSessionToken(r.Context(), sqlcore.GetSessionTokenParams{
		ID:     id,
		UserID: userID,
	})
}

func (s *sessions) UpdateToken(r *http.Request, id uuid.UUID, token string) error {
	return s.queries.UpdateSessionToken(r.Context(), sqlcore.UpdateSessionTokenParams{
		ID:    id,
		Token: token,
	})
}

func (s *sessions) DeleteAll(r *http.Request, id uuid.UUID) error {
	return s.queries.DeleteUserSessions(r.Context(), id)
}

func (s *sessions) Delete(r *http.Request, id uuid.UUID, userID uuid.UUID) error {
	return s.queries.DeleteUserSession(r.Context(), sqlcore.DeleteUserSessionParams{
		ID:     id,
		UserID: userID,
	})
}
