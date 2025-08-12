package response

import (
	sessionProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func Session(session models.Session) *sessionProto.Session {
	return &sessionProto.Session{
		Id:        session.ID.String(),
		UserId:    session.UserID.String(),
		Client:    session.Client,
		Ip:        session.IP,
		CreatedAt: session.CreatedAt.String(),
		LastUsed:  session.LastUsed.String(),
	}
}
