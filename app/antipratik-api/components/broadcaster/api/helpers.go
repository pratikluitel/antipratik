package api

import (
	"net/http"
	"strconv"
)

// ── Helpers ───────────────────────────────────────────────────────────────────

func parseBroadcastID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	raw := r.PathValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}
