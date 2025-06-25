package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/auth"
)

func SessionList(sessions []models.Session) *auth.SessionsListResponse {

	list := make([]*auth.SessionResponse, len(sessions))
	for i, session := range sessions {
		list[i] = &auth.SessionResponse{
			Id:        session.ID.String(),
			UserId:    session.UserID.String(),
			Client:    session.Client,
			Ip:        session.IP,
			CreatedAt: session.CreatedAt.String(),
			LastUsed:  session.LastUsed.String(),
		}
	}

	return &auth.SessionsListResponse{
		Sessions: list,
	}
}
