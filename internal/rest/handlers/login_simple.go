package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/rest/responses"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/jsonkit"
	"github.com/chains-lab/gatekit/roles"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (h *Handlers) LoginSimple(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New()
	log := h.log.WithField("request_id", requestID)

	if !h.cfg.Server.TestMode {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusForbidden,
			Title:  "Test mode is off",
			Detail: "Test mode is off",
		})...)
	}

	type emailReq struct {
		Email string `json:"email"`
		Role  string `json:"role"`
		//Sub   *string `json:"sub,omitempty"`
	}
	var req emailReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
		"role":  validation.Validate(req.Role, validation.NilOrNotEmpty),
	}
	if errs.Filter() != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  errs.Filter(),
		})...)
		return
	}

	role, err := roles.ParseRole(req.Role)
	if err != nil {
		log.WithError(err).Error("failed to parse role")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Title:  "Invalid role",
			Detail: "The provided role is invalid",
		})...)
		return
	}

	session, tokensPair, appErr := h.app.Login(r.Context(), req.Email, role, r.Header.Get("User-Agent"))
	if appErr != nil {
		log.WithError(appErr.Unwrap()).Error("error getting session")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		return
	}

	log.Debugf("got session: %+v", session)
	httpkit.Render(w, responses.TokensPair(tokensPair.Access, tokensPair.Refresh))
}
