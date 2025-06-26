package responses

import (
	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func Session(session models.Session) *svc.SessionResponse {
	return &svc.SessionResponse{
		Id:        session.ID.String(),
		UserId:    session.UserID.String(),
		Client:    session.Client,
		Ip:        session.IP,
		CreatedAt: session.CreatedAt.String(),
		LastUsed:  session.LastUsed.String(),
	}
}
