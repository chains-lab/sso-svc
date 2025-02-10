package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/cifractx"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/resources"
	"github.com/recovery-flow/tokens"
	"github.com/sirupsen/logrus"
)

func SessionDelete(w http.ResponseWriter, r *http.Request) {
	Server, err := cifractx.GetValue[*config.Server](r.Context(), config.SERVER)
	if err != nil {
		logrus.Errorf("Failed to retrieve service configuration %s", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}
	log := Server.Logger

	sessionForDeleteId, err := uuid.Parse(chi.URLParam(r, "session_id"))
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	sessionID, ok := r.Context().Value(tokens.DeviceIDKey).(uuid.UUID)
	if !ok {
		log.Warn("SessionID not found in context")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	userID, ok := r.Context().Value(tokens.UserIDKey).(uuid.UUID)
	if !ok {
		log.Warn("UserID not found in context")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if sessionID == sessionForDeleteId {
		log.Debugf("Session can't be current")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("session can't be current"))...)
		return
	}

	err = Server.SqlDB.Sessions.Delete(r, sessionForDeleteId, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpkit.RenderErr(w, problems.NotFound())
			return
		}
		log.Errorf("Failed to delete device: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	err = Server.TokenManager.Bin.Add(userID.String(), sessionForDeleteId.String())
	if err != nil {
		log.Errorf("Failed to add token to bin: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	sessions, err := Server.SqlDB.Sessions.GetSessions(r, userID)
	if err != nil {
		log.Errorf("Failed to retrieve user sessions: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	var userSessions []resources.Session
	for _, device := range sessions {
		userSessions = append(userSessions, NewSession(r, device.ID.String()))
	}

	httpkit.Render(w, NewSessionList(userSessions))
}

func NewSessionList(sessions []resources.Session) resources.UserSessions {
	return resources.UserSessions{
		Data: resources.UserSessionsData{
			Type: resources.UserSessionsType,
			Attributes: resources.UserSessionsDataAttributes{
				Devices: sessions,
			},
		},
	}
}

func NewSession(r *http.Request, Id string) resources.Session {
	return resources.Session{
		Id:        Id,
		Client:    httpkit.GetUserAgent(r),
		IpFirst:   httpkit.GetClientIP(r),
		IpLast:    httpkit.GetClientIP(r),
		LastUsed:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
}
