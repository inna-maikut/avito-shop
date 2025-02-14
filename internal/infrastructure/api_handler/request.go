package api_handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func Parse[T any](r *http.Request, w http.ResponseWriter, t *T) (ok bool) {
	bodyBytes, err := io.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()
	if err != nil {
		BadRequest(w, "could not read request body")
		return false
	}

	err = json.Unmarshal(bodyBytes, t)
	if err != nil {
		BadRequest(w, "could not bind request body")
		return false
	}

	return true
}
