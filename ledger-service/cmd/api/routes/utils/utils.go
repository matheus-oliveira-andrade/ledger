package utils

import (
	"encoding/json"
	"net/http"
)

func RenderErrorJsonResponse(w http.ResponseWriter, err error) {
	errObj := struct {
		ErrorMessage string `json:"errorMessage"`
	}{
		ErrorMessage: err.Error(),
	}

	RenderJsonResponse(w, errObj, http.StatusBadRequest)
}

func RenderJsonResponse(w http.ResponseWriter, obj any, httpStatus int) {
	w.WriteHeader(httpStatus)

	if obj == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		panic(err)
	}
}
