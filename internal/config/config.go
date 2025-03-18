package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Host                       string        `json:"host"`
	Port                       string        `json:"port"`
	ReadTimeoutInMilliseconds  time.Duration `jsom:"readTimeoutInMilliseconds"`
	WriteTimeoutInMilliseconds time.Duration `jsom:"WriteTimeoutInMilliseconds"`

	DbOptions struct {
		DriverName string `json:"driver"`
		DbPath     string `json:"path"`
	}

	EmailOptions struct {
		Host string `json:"emailHost"`
		Port string `json:"emailPort"`
	} `json:"emailOptions"`

	ListSenderOptions struct {
		Delay time.Duration `json:"delayInSeconds"`
	} `json:"listSenderOptions"`
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
