package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/golang-jwt/jwt/v5"
)

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

func readJWT(rawTok string, key []byte) (jwt.MapClaims, error) {
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
