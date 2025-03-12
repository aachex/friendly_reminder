package controller

import (
	"encoding/json"
	"io"
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
