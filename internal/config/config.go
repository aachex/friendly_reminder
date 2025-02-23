package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port          int    `json:"port"`
	Email         string `json:"email"`
	EmailPassword string `json:"emailPassword"`
	EmailHost     string `json:"emailHost"`
	EmailPort     int    `json:"emailPort"`
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
