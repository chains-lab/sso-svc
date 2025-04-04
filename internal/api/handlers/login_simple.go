package handlers

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/sso-oauth/internal/api/responses"
	"github.com/hs-zavet/sso-oauth/internal/app"
)

func (h *Handler) LoginSimple(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.Server.TestMode {
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
		"role":  validation.Validate(req.Role, validation.NilOrNotEmpty),
		"sub":   validation.Validate(req.Sub, validation.NilOrNotEmpty),
	}
	if errs.Filter() != nil {
		httpkit.RenderErr(w, problems.BadRequest(errs.Filter())...)
		return
	}

	res, err := h.app.Login(r.Context(), app.LoginRequest{
		Email: req.Email,
	})
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.TokensPair(res.Access, res.Refresh))
}
