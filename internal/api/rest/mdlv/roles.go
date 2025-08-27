package mdlv

import (
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/httplab"
	"github.com/google/jsonapi"
)

func AccessGrant(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(meta.UserCtxKey).(meta.UserData)
			if !ok {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "UNAUTHORIZED",
					Detail: "User not authenticated",
				})
				return
			}

			roleAllowed := false
			for _, role := range roles {
				if user.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusForbidden),
					Code:   "USER_ROLE_NOT_ALLOWED",
					Detail: "User role not allowed",
				})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
