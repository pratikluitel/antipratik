package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// ── Shared helpers ────────────────────────────────────────────────────────────

func streamFile(w http.ResponseWriter, r *http.Request, body io.ReadSeekCloser, ct string) {
	defer body.Close()
	w.Header().Set("Content-Type", ct)
	http.ServeContent(w, r, "", time.Time{}, body)
}

// formTags returns nil when neither "tags" nor "tags[]" key is present in the form
// (keep existing), or the slice of values (possibly empty = clear all) when either
// key is present. Browsers/FormData append array fields as "tags[]"; plain JSON forms
// use "tags". Values may be comma-separated (e.g. "tag1,tag2") and are split and trimmed.
func formTags(r *http.Request) []string {
	key := "tags"
	if _, ok := r.Form["tags[]"]; ok {
		key = "tags[]"
	} else if _, ok := r.Form["tags"]; !ok {
		// For multipart requests, absence of tags[] means clear all tags.
		// Non-multipart callers (JSON bodies) preserve existing tags when the key is absent.
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			return []string{}
		}
		return nil
	}

	var tags []string
	for _, v := range r.Form[key] {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				tags = append(tags, part)
			}
		}
	}
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		return []string{}
	}
	return tags
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Encoding error here means our own struct failed to serialise —
	// the response is already partially written and unrecoverable, so discard silently.
	_ = json.NewEncoder(w).Encode(v)
}
