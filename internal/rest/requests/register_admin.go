package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/restkit/roles"
	"github.com/chains-lab/sso-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func newDecodeError(what string, err error) error {
	return validation.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}

func RegisterAdmin(r *http.Request) (req resources.RegisterAdmin, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.RegisterAdminType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),

		"data/attributes/email": validation.Validate(
			req.Data.Attributes.Email, validation.Required, validation.Length(5, 255), is.Email),

		"data/attributes/role": validation.Validate(
			req.Data.Attributes.Role, validation.Required, validation.In(roles.GetAllRoles())),
	}

	//TODO
	//if err = passmanager.New().ReliabilityCheck(req.Data.Attributes.Password); err != nil {
	//	errs["data/attributes/password"] = err
	//}

	return req, errs.Filter()
}
