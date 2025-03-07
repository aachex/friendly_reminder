package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Email struct {
		Host string `json:"emailHost"`
		Port string `json:"emailPort"`
	} `json:"email"`
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
