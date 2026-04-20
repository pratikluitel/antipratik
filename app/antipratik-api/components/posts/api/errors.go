package api

import (
	"net/http"

	commonerrors "github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
)

// handleLogicError writes a 400 Bad Request for validation errors and a 500
// Internal Server Error for all other failures. The operation name is included
// in the internal log message but never exposed to the client.
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
