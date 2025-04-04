package handlers

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/tokens/identity"
	"github.com/pkg/errors"
)

func LoginSimple(w http.ResponseWriter, r *http.Request) {
	if !Config(r).Server.TestMode {
		Log(r).Warn("Test mode is off")
		httpkit.RenderErr(w, problems.Forbidden("Test mode is off"))
	}

	type emailReq struct {
		Email string  `json:"email"`
		Role  string  `json:"role"`
		Sub   *string `json:"sub,omitempty"`
	}
	var req emailReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
		"role":  validation.Validate(req.Role, validation.Required),
		"sub":   validation.Validate(req.Sub, validation.NilOrNotEmpty),
	}
	if errs.Filter() != nil {
		Log(r).WithError(errs.Filter()).Error("Failed to parse email")
		httpkit.RenderErr(w, problems.BadRequest(errs.Filter())...)
		return
	}

	role, err := identity.ParseIdentityType(req.Role)
	if err != nil {
		Log(r).WithError(err).Error("Invalid role")
		httpkit.RenderErr(w, problems.BadRequest(errors.New("invalid role"))...)
		return
	}

	var sub *uuid.UUID
	if req.Sub != nil {
		id, err := uuid.Parse(*req.Sub)
		if err != nil {
			Log(r).WithError(err).Error("Invalid sub")
			httpkit.RenderErr(w, problems.BadRequest(errors.New("invalid sub"))...)
			return
		}
		sub = &id
	}

	tokenAccess, tokenRefresh, err := Domain(r).Login(r.Context(), role, sub, req.Email, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		Log(r).WithError(err).Error("Failed to login")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	Log(r).WithField("tokenAccess", tokenAccess).Debugf("Successfully logged in")

	httpkit.Render(w, responses.TokensPair(*tokenAccess, *tokenRefresh))
}
