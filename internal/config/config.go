package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

func NewConfig(path string) *Config {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var conf Config
	json.Unmarshal(fileContent, &conf)
	return &conf
}
