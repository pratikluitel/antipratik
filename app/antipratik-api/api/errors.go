package api

import (
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/logic"
)

// handleLogicError writes a 400 Bad Request for validation errors and a 500
// Internal Server Error for all other failures. The operation name is included
// in the internal log message but never exposed to the client.
func handleLogicError(w http.ResponseWriter, op string, err error) {
	if logic.IsValidationError(err) {
		writeError(w, http.StatusBadRequest, err.Error())
	} else {
		log.Printf("%s error: %v", op, err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
