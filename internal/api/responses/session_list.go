package responses

import (
	svc "github.com/chains-lab/proto-storage/gen/go/sso"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func SessionList(sessions []models.Session) *svc.SessionsListResponse {

	list := make([]*svc.SessionResponse, len(sessions))
	for i, session := range sessions {
		list[i] = &svc.SessionResponse{
			Id:        session.ID.String(),
			UserId:    session.UserID.String(),
			Client:    session.Client,
			Ip:        session.IP,
			CreatedAt: session.CreatedAt.String(),
			LastUsed:  session.LastUsed.String(),
		}
	}

	return &svc.SessionsListResponse{
		Sessions: list,
	}
}
