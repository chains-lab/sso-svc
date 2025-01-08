package requests

import (
	"encoding/json"
	"net/http"

	"github.com/cifra-city/sso-oauth/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func NewRefresh(r *http.Request) (req resources.RefreshToken, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.RefreshType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
