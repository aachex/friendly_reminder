package jwtutil

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Parse преобразует jwt токен в структуру [jwt.Token].
func Parse(rawJwt string, key []byte) (*jwt.Token, error) {
	tok, err := jwt.Parse(rawJwt, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("expected signing method HMAC, got %s", t.Method.Alg())
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !tok.Valid {
		return nil, errors.New("token is invalid")
	}

	return tok, err
}

// GetClaims получает полезные данные из jwt токена.
func GetClaims(rawJwt string, key []byte) (jwt.MapClaims, error) {
	tok, err := Parse(rawJwt, key)
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid jwt claims")
	}

	return claims, nil
}

// FromHeader возвращает закодированный токен из заголовка запроса.
func FromHeader(h http.Header) string {
	rawTok := h.Get("Authorization")
	if len(rawTok) < 8 {
		return ""
	}
	return rawTok[7:] // Отрезаем часть 'Bearer '
}
