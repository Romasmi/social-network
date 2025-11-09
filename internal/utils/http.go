package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ErrorInvalidRequestBody(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("invalid request body, %s", err.Error()), http.StatusBadRequest)
}

func JsonResponse(w http.ResponseWriter, output interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(output)
	if err != nil {
		fmt.Printf("error while encoding response: %v\n", err)
		JsonError(w, fmt.Errorf("internal error"))
	}
}

func JsonError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func JsonErrorNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
}
