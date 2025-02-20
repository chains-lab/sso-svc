package handlers

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/comtools/jsonkit"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
	"github.com/recovery-flow/tokens/identity"
)

func (a *App) LoginSimple(w http.ResponseWriter, r *http.Request) {
	if !a.Config.Server.TestMode {
		a.Log.Warn("Test mode is off")
		httpkit.RenderErr(w, problems.Forbidden("Test mode is off"))
	}

	type emailReq struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	var req emailReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
		"role":  validation.Validate(req.Role, validation.Required),
	}
	if errs.Filter() != nil {
		a.Log.WithError(errs.Filter()).Error("Failed to parse email")
		httpkit.RenderErr(w, problems.BadRequest(errs.Filter())...)
		return
	}

	role, err := identity.ParseIdentityType(req.Role)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(errors.New("invalid role"))...)
		return
	}

	tokenAccess, tokenRefresh, err := a.Domain.SessionLogin(r.Context(), role, req.Email, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
