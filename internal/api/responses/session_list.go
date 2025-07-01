package responses

import (
	svc "github.com/chains-lab/proto-storage/gen/go/svc/sso"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func SessionList(sessions []models.Session) *svc.SessionsList {

	list := make([]*svc.Session, len(sessions))
	for i, session := range sessions {
		list[i] = &svc.Session{
			Id:        session.ID.String(),
			UserId:    session.UserID.String(),
			Client:    session.Client,
			Ip:        session.IP,
			CreatedAt: session.CreatedAt.String(),
			LastUsed:  session.LastUsed.String(),
		}
	}

	return &svc.SessionsList{
		Sessions: list,
	}
}
