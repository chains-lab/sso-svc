package meta

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/auth"
)

type ctxKey int

const (
	UserCtxKey ctxKey = iota
	IpCtxKey
	UserAgentCtxKey
	ClientTxKey
)

func User(ctx context.Context) (auth.UserData, error) {
	if ctx == nil {
		return auth.UserData{}, fmt.Errorf("mising context")
	}

	userData, ok := ctx.Value(UserCtxKey).(auth.UserData)
	if !ok {
		return auth.UserData{}, fmt.Errorf("mising context")
	}

	return userData, nil
}

//func Ip(r *http.Request) string {
//	return r.Context().Value(IpCtxKey).(string)
//}
//
//func UserAgent(r *http.Request) string {
//	return r.Context().Value(UserAgentCtxKey).(string)
//}
//
//func Client(r *http.Request) string {
//	return r.Context().Value(ClientTxKey).(string)
//}
