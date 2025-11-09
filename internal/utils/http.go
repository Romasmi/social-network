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
	err := json.NewEncoder(w).Encode(output)
	w.Header().Set("Content-Type", "application/json")
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
