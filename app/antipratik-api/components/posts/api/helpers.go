package api

import (
	"net/http"
	"strings"

	"github.com/pratikluitel/antipratik/common/errors"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/common/requests"
)

// handleLogicError writes a 400 Bad Request for validation errors and a 500
// Internal Server Error for all other failures. The operation name is included
// in the internal log message but never exposed to the client.
func handleLogicError(w http.ResponseWriter, log logging.Logger, op string, err error) {
	if errors.Is(err) {
		requests.WriteError(w, http.StatusBadRequest, err.Error())
	} else {
		log.Error(op+" error", "err", err)
		requests.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
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
