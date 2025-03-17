package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errReadingBody = errors.New("error reading request body")
)

var key = []byte(os.Getenv("SECRET_STR"))

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

func readJWT(rawTok string) (jwt.MapClaims, error) {
	tok, err := jwt.Parse(rawTok, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid jwt claims")
}

func getRawJwtFromHeader(h http.Header) string {
	rawTok := h.Get("Authorization")
	if len(rawTok) > 7 {
		return rawTok[7:] // Отрезаем часть 'Bearer '
	}
	return ""
}
