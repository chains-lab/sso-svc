package responses

import (
	"github.com/chains-lab/sso-svc/internal/domain/entity"
	"github.com/chains-lab/sso-svc/resources"
)

func AccountSession(m entity.Session) resources.AccountSession {
	resp := resources.AccountSession{
		Data: resources.AccountSessionData{
			Id:   m.ID,
			Type: resources.AccountSessionType,
			Attributes: resources.AccountSessionAttributes{
				AccountId: m.AccountID,
				CreatedAt: m.CreatedAt,
				LastUsed:  m.LastUsed,
			},
		},
	}

	return resp
}

func AccountSessionsCollection(ms entity.SessionsCollection) resources.AccountSessionsCollection {
	items := make([]resources.AccountSessionData, 0, len(ms.Data))

	for _, s := range ms.Data {
		items = append(items, AccountSession(s).Data)
	}

	return resources.AccountSessionsCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
