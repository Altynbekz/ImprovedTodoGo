package utils

import (
	"encoding/json"
	"net/http"
)

func JsonResponseError(w http.ResponseWriter, code int, msg string) {
	JsonResponse(w, code, msg, nil)
}

func JsonResponse(w http.ResponseWriter, code int, comment string, payload any) {
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Code    int    `json:"code"`
		Comment string `json:"comment"`
		Data    any    `json:"data,omitempty"`
	}{
		Code:    code,
		Comment: comment,
		Data:    payload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":500,"comment":"Something went wrong"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithoutBody(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
