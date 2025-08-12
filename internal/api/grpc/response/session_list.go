package response

import (
	pagProto "github.com/chains-lab/sso-proto/gen/go/common/pagination"
	sessionProto "github.com/chains-lab/sso-proto/gen/go/svc/session"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/internal/pagination"
)

func SessionList(sessions []models.Session, response pagination.Response) *sessionProto.SessionsList {
	list := make([]*sessionProto.Session, len(sessions))
	for i, session := range sessions {
		list[i] = &sessionProto.Session{
			Id:        session.ID.String(),
			UserId:    session.UserID.String(),
			Client:    session.Client,
			Ip:        session.IP,
			CreatedAt: session.CreatedAt.String(),
			LastUsed:  session.LastUsed.String(),
		}
	}

	return &sessionProto.SessionsList{
		Sessions: list,
		Pagination: &pagProto.Response{
			Page:  response.Page,
			Size:  response.Size,
			Total: response.Total,
		},
	}
}
