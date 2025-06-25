package responses

import (
	"github.com/chains-lab/chains-auth/internal/app/models"
	"github.com/chains-lab/proto-storage/gen/go/auth"
)

func Session(session models.Session) *auth.SessionResponse {
	return &auth.SessionResponse{
		Id:        session.ID.String(),
		UserId:    session.UserID.String(),
		Client:    session.Client,
		Ip:        session.IP,
		CreatedAt: session.CreatedAt.String(),
		LastUsed:  session.LastUsed.String(),
	}
}
