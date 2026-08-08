package util

import (
	"net/http"
)

const (
	contentTypeKey = "Content-Type"
	jsonMediaType  = "application/json"
)

func WriteJSON(w http.ResponseWriter, dto any, responseStatus int) error {
	bytes, err := ToJSON(dto)
	if err != nil {
		return err
	}
	w.Header().Set(contentTypeKey, jsonMediaType)
	w.WriteHeader(responseStatus)
	_, err = w.Write(bytes)
	return err
}
