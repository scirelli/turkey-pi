package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/scirelli/turkey-pi/internal/app/server"
	"github.com/scirelli/turkey-pi/pkg/log"
)

const (
	KEYBOARD_DEFAULT_FILE string = "/dev/hidg0"
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

	Defaults(&config)

	return &config, nil
}

func Defaults(config *AppConfig) *AppConfig {
	var logger = log.New("AppConfig", log.GetLevel(config.LogLevel))

	if config.LogLevel == "" {
		config.Server.LogLevel = log.DEFAULT_LOG_LEVEL
		logger.Infof("Defaulting server log level to '%s'", config.Server.LogLevel)
	} else {
		config.Server.LogLevel = log.GetLevel(config.LogLevel)
		logger.Infof("Setting server log level to '%s'", config.Server.LogLevel)
	}

	if config.KeyboardFile == "" {
		config.KeyboardFile = KEYBOARD_DEFAULT_FILE
		logger.Infof("Defaulting keyboard file to '%s'", config.KeyboardFile)
	}

	config.Server.Debug = config.Debug

	server.Defaults(&config.Server)

	return config
}

//AppConfig configuration data for entire application.
type AppConfig struct {
	Debug              bool              `json:"debug"`
	LogLevel           string            `json:"logLevel"`
	KeyboardFile       string            `json:"keyboardFile"`
	CharacterToKeyFile string            `json:"characterToKeyFile,omitempty"`
	CharacterToKeyMap  map[string]string `json:"characterToKeyMap"`
	Server             server.Config     `json: "server,omitempty"`
}
