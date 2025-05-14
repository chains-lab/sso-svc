package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/chains-auth/internal/app"
	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/jsonkit"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (h *Handler) LoginSimple(w http.ResponseWriter, r *http.Request) {
	if !h.cfg.Server.TestMode {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusForbidden,
			Title:  "Test mode is off",
			Detail: "Test mode is off",
		})...)
	}

	type emailReq struct {
		Email string `json:"email"`
		//Role  string  `json:"role"`
		//Sub   *string `json:"sub,omitempty"`
	}
	var req emailReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"email": validation.Validate(req.Email, validation.Required),
		//"role":  validation.Validate(req.Role, validation.NilOrNotEmpty),
		//"sub":   validation.Validate(req.Sub, validation.NilOrNotEmpty),
	}
	if errs.Filter() != nil {
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  errs.Filter(),
		})...)
		return
	}

	res, err := h.app.Login(r.Context(), app.LoginRequest{
		Email: req.Email,
	})
	if err != nil {
		h.log.WithError(err).Error("error getting session")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusInternalServerError,
		})...)
		return
	}

	h.log.Infof("got session: %+v", res)
	httpkit.Render(w, responses.TokensPair(res.Access, res.Refresh))
}
