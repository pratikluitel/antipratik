package requests

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, Error{Message: msg})
}
