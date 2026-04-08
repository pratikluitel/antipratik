package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// ── Shared helpers ────────────────────────────────────────────────────────────

func streamFile(w http.ResponseWriter, r *http.Request, body io.ReadCloser, ct string) {
	defer body.Close()
	w.Header().Set("Content-Type", ct)
	if _, err := io.Copy(w, body); err != nil && r.Context().Err() == nil {
		log.Printf("streamFile copy error: %v", err)
	}
}

// formTags returns nil when the "tags" key is absent from the form (keep existing),
// or the slice of values (possibly empty = clear all) when the key is present.
// Values may be comma-separated (e.g. "tag1,tag2") and are split and trimmed.
func formTags(r *http.Request) []string {
	key := "tags"
	if _, ok := r.Form[key]; !ok {
		return nil
	} else if _, ok := r.Form["tags[]"]; ok {
		key = "tags[]"
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
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
