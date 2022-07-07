package server

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/scirelli/turkey-pi/pkg/log"
)

const (
	DEFAULT_PORT uint = 8282
)

//Load a config file.
func Load(fileName string) (*Config, error) {
	var config Config

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

func Defaults(config *Config) *Config {
	var logger = log.New("ServerConfig", config.LogLevel)

	if config.ContentPath == "" {
		config.ContentPath = "."
	}
	if config.Port == 0 {
		config.Port = DEFAULT_PORT
	}
	if config.ServerUrl == "" {
		config.ServerUrl = "http://localhost:8282"
		logger.Infof("Defaulting ServerUrl to '%s'\n", config.ServerUrl)
	} else {
		base, err := url.Parse(config.ServerUrl)
		if err != nil {
			logger.Panic(err)
		}
		config.ServerUrl = base.String()
	}
	if config.UiUrl == "" {
		config.UiUrl = "http://localhost"
		logger.Infof("Defaulting UiUrl to '%s'\n", config.UiUrl)
	}

	return config
}

type Config struct {
	Port    uint   `json:"port"`
	Address string `json:"address"`

	ContentPath string `json:"contentPath"`

	ServerUrl string `json:"serverUrl"`
	UiUrl     string `json:"uiUrl"`

	Debug    bool         `json:"-"`
	LogLevel log.LogLevel `json:"-"`
}
