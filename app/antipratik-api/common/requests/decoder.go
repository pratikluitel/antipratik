package requests

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		slog.Error("failed to decode JSON body", "error", err)
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return err
	}
	return nil
}
