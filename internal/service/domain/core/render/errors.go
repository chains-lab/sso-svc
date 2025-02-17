package render

import (
	"database/sql"
	"net/http"

	"github.com/pkg/errors"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/sirupsen/logrus"
)

func RenderSelectErr(w http.ResponseWriter, log *logrus.Logger, err error, message string) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		httpkit.RenderErr(w, problems.Unauthorized())
	default:
		log.Errorf("%s: %v", message, err)
		httpkit.RenderErr(w, problems.InternalError())
	}
}
