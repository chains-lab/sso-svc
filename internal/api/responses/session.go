package responses

import (
	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func Session(session models.Session) *svc.Session {
	return &svc.Session{
		Id:        session.ID.String(),
		UserId:    session.UserID.String(),
		Client:    session.Client,
		Ip:        session.IP,
		CreatedAt: session.CreatedAt.String(),
		LastUsed:  session.LastUsed.String(),
	}
}
