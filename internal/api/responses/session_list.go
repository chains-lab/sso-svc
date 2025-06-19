package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/sso"
)

func SessionList(sessions []models.Session) *sso.SessionsListResponse {

	list := make([]*sso.SessionResponse, len(sessions))
	for i, session := range sessions {
		list[i] = &sso.SessionResponse{
			Id:        session.ID.String(),
			UserId:    session.UserID.String(),
			Client:    session.Client,
			Ip:        session.IP,
			CreatedAt: session.CreatedAt.String(),
			LastUsed:  session.LastUsed.String(),
		}
	}

	return &sso.SessionsListResponse{
		Sessions: list,
	}
}
