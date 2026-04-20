package api

import (
	"net/http"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
)

func handleLogicError(w http.ResponseWriter, log logging.Logger, op string, err error) {
	if commonerrors.Is(err) {
		writeError(w, http.StatusBadRequest, err.Error())
	} else {
		log.Error(op+" error", "err", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
