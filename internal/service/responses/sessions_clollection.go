package responses

import (
	"github.com/recovery-flow/sso-oauth/internal/data/sql/repositories/sqlcore"
	"github.com/recovery-flow/sso-oauth/resources"
)

func SessionCollection(sessions []sqlcore.Session) resources.SessionsCollection {
	var data []resources.SessionData
	for _, session := range sessions {
		data = append(data, Session(session).Data)
	}
	return resources.SessionsCollection{
		Data: resources.SessionsCollectionData{
			Type: resources.UserSessionsType,
			Attributes: resources.SessionsCollectionDataAttributes{
				Sessions: data,
			},
		},
	}
}
