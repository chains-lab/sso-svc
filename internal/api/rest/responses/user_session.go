package responses

import (
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/sso-svc/internal/app/models"
	"github.com/chains-lab/sso-svc/resources"
)

func UserSession(m models.Session) resources.UserSession {
	resp := resources.UserSession{
		Data: resources.UserSessionData{
			Id:   m.ID.String(),
			Type: resources.UserSessionType,
			Attributes: resources.UserSessionAttributes{
				UserId:    m.UserID.String(),
				CreatedAt: m.CreatedAt,
				LastUsed:  m.LastUsed,
			},
		},
	}

	return resp
}

func UserSessionsCollection(ms []models.Session, pag pagi.Response) resources.UserSessionsCollection {
	items := make([]resources.UserSessionData, 0, len(ms))

	for _, s := range ms {
		items = append(items, UserSession(s).Data)
	}

	return resources.UserSessionsCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}
}
