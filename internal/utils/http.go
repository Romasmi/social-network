package utils

import "net/http"

func ErrorInvalidRequestBody(w http.ResponseWriter) {
	http.Error(w, "invalid request body", http.StatusBadRequest)
}
