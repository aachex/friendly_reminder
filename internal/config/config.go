package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port      string `json:"port"`
	EmailHost string `json:"emailHost"`
	EmailPort string `json:"emailPort"`
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
