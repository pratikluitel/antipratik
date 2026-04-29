package requests

import (
	"encoding/json"
	"net/http"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return err
	}
	return nil
}
