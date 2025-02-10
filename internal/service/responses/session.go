package responses

import (
	"github.com/recovery-flow/sso-oauth/internal/data/sql/repositories/sqlcore"
	"github.com/recovery-flow/sso-oauth/resources"
)

func Session(session sqlcore.Session) resources.Session {
	return resources.Session{
		Data: resources.SessionData{
			Type: resources.UserSessionType,
			Id:   session.ID.String(),
			Attributes: resources.SessionAttributes{
				UserId:    session.UserID.String(),
				Client:    session.Client,
				Ip:        session.Ip,
				CreatedAt: session.CreatedAt,
				LastUsed:  session.LastUsed,
			},
		},
	}
}
