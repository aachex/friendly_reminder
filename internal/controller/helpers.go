package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

var (
	errReadingBody = errors.New("error reading request body")
)

var jwtKey = []byte(os.Getenv("SECRET_STR"))

func readBody[T any](body io.ReadCloser) (*T, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var t T
	if err = json.Unmarshal(bodyBytes, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func writeJson[T any](w http.ResponseWriter, obj T) {
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getRawJwtFromHeader(h http.Header) string {
	rawTok := h.Get("Authorization")
	if len(rawTok) < 8 {
		return ""
	}
	return rawTok[7:] // Отрезаем часть 'Bearer '
}
