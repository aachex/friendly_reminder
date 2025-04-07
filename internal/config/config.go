package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	Prefix       string        `json:"apiPrefix"`
	ReadTimeout  time.Duration `jsom:"readTimeout"`
	WriteTimeout time.Duration `jsom:"writeTimeout"`

	Sqlite3 database `json:"sqlite3"`

	EmailOptions struct {
		Host string `json:"emailHost"`
		Port string `json:"emailPort"`
	} `json:"emailOptions"`

	ListSenderOptions struct {
		Delay time.Duration `json:"delay"`
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

type database struct {
	DriverName string `json:"driver"`
	ConnStr    string `json:"connStr"`
}
