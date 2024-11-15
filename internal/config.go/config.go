package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Sql     SqlConfig     `json:"sqlConfig"`
	Server  ServerConfig  `json:"serverConfig"`
	Checker CheckerConfig `json:"checkerConfig"`
}

type CheckerConfig struct {
	Delay int64 `json:"delay"`
}

type SqlConfig struct {
	Storagepath string `json:"storagepath"`
}

type ServerConfig struct {
	Port              string `json:"port"`
	ReadTimeout       int    `json:"readTimeout"`
	WriteTimeout      int    `json:"writeTimeout"`
	IdleTimeout       int    `json:"idleTimeout"`
	ReadHeaderTimeout int    `json:"readHeaderTimeout"`
}

func MustLoad(configPath string) Config {
	var cfg Config

	file, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalln("err on opening config file, err: ", err)
	}

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatalln("err on decoding config file, err: ", err)
	}

	return cfg
}
