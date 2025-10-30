package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/sso-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdatePassword(r *http.Request) (req resources.UpdatePassword, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UpdatePasswordType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if req.Data.Attributes.NewPassword != req.Data.Attributes.ConfirmPassword {
		errs["data/attributes/confirm_password"] = fmt.Errorf("must match password")
	}

	//TODO
	//if err = passmanager.New().ReliabilityCheck(req.Data.Attributes.NewPassword); err != nil {
	//	errs["data/attributes/new_password"] = err
	//}

	return req, errs.Filter()
}
