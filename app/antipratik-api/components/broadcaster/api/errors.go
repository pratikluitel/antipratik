package api

import (
	"net/http"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
)

func handleLogicError(w http.ResponseWriter, log logging.Logger, op string, err error) {
	if commonerrors.Is(err) {
		requests.WriteError(w, http.StatusBadRequest, err.Error())
	} else {
		log.Error(op+" error", "err", err)
		requests.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
