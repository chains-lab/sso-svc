package responses

import (
	"github.com/chains-lab/sso-svc/internal/models"
	"github.com/chains-lab/sso-svc/resources"
)

func UserSession(m models.Session) resources.UserSession {
	resp := resources.UserSession{
		Data: resources.UserSessionData{
			Id:   m.ID,
			Type: resources.UserSessionType,
			Attributes: resources.UserSessionAttributes{
				UserId:    m.UserID,
				CreatedAt: m.CreatedAt,
				LastUsed:  m.LastUsed,
			},
		},
	}

	return resp
}

func UserSessionsCollection(ms models.SessionsCollection) resources.UserSessionsCollection {
	items := make([]resources.UserSessionData, 0, len(ms.Data))

	for _, s := range ms.Data {
		items = append(items, UserSession(s).Data)
	}

	return resources.UserSessionsCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
