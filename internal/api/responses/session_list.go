package responses

import (
	"github.com/chains-lab/sso-proto/gen/go/svc"
	"github.com/chains-lab/sso-svc/internal/app/models"
)

func SessionList(sessions []models.Session, page uint64) *svc.SessionsList {

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
		Pagination: &svc.PaginationResponse{
			Page:  page,
			Total: uint64(len(sessions)),
		},
	}
}
