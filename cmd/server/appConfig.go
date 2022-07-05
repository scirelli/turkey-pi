package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/scirelli/turkey-pi/internal/app/server"
	"github.com/scirelli/turkey-pi/pkg/log"
)

//LoadConfig a config file.
func LoadConfig(fileName string) (*AppConfig, error) {
	var config AppConfig

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return &config, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return &config, err
	}

	json.Unmarshal(byteValue, &config)

	if config.LogLevel == "" {
		config.Server.LogLevel = log.DEFAULT_LOG_LEVEL
	} else {
		config.Server.LogLevel = log.GetLevel(config.LogLevel)
	}

	server.Defaults(&config.Server)

	return &config, nil
}

//AppConfig configuration data for entire application.
type AppConfig struct {
	Debug              bool              `json:"debug"`
	LogLevel           string            `json:"logLevel"`
	CharacterToKeyFile string            `json:"characterToKeyFile,omitempty"`
	CharacterToKeyMap  map[string]string `json:"characterToKeyMap"`
	Server             server.Config     `json: "server,omitempty"`
}
