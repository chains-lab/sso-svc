package mdlv

import (
	"context"
	"net/http"
	"strings"

	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/httplab"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
)

const (
	hAuthorization        = "Authorization"
	hServiceAuthorization = "X-Service-Authorization" // отдельный заголовок для m2m
	bearerPrefix          = "bearer "
)

func AuthMdl(skUser string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(hAuthorization)
			if authHeader == "" {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusNotFound),
					Code:   "MISSING_AUTHORIZATION_HEADER",
					Detail: "Missing Authorization header",
				})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "MISSING_AUTHORIZATION_HEADER",
					Detail: "Missing Authorization header",
				})
				return
			}

			tokenString := parts[1]

			userData, err := auth.VerifyUserJWT(r.Context(), tokenString, skUser)
			if err != nil {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "TOKEN_VALIDATION_FAILED",
					Detail: "Token validation failed",
				})
				return
			}

			userID, err := uuid.Parse(userData.Subject)
			if err != nil {
				httplab.RenderErr(w, &jsonapi.ErrorObject{
					Status: http.StatusText(http.StatusUnauthorized),
					Code:   "INVALID_USER_ID",
					Detail: "Invalid user ID in token",
				})
				return
			}

			ctx = context.WithValue(ctx, meta.UserCtxKey, meta.UserData{
				ID:        userID,
				SessionID: userData.Session,
				Role:      userData.Role,
				Verified:  userData.Verified,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
