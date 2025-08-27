package mdlv

import (
	"net/http"
	"strings"

	"github.com/chains-lab/distributors-svc/internal/config/constant"
	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/httplab"
	"github.com/google/jsonapi"
)

func ServiceAuthMdl(skService string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			h := r.Header.Get(hServiceAuthorization)
			if h == "" {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "MISSING_SERVICE_AUTH_HEADER",
					Detail: "Missing X-Service-Authorization header",
				})
				return
			}

			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "INVALID_SERVICE_AUTH_SCHEME",
					Detail: "X-Service-Authorization must be Bearer <token>",
				})
				return
			}

			svcToken := parts[1]

			svcData, err := auth.VerifyServiceJWT(ctx, svcToken, skService)
			if err != nil {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "SERVICE_TOKEN_VALIDATION_FAILED",
					Detail: "Service token validation failed",
				})
				return
			}

			access := false
			for _, v := range svcData.Audience {
				if v == constant.ServiceName {
					access = true
					break
				}
			}
			if !access {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusForbidden),
					Code:   "FORBIDDEN",
					Detail: "Service does not have access",
				})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
