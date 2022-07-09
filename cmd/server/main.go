package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/scirelli/turkey-pi/internal/app/server"
	"github.com/scirelli/turkey-pi/pkg/keyboard"
	"github.com/scirelli/turkey-pi/pkg/log"
)

func main() {
	var logger = log.New("Main", log.DEFAULT_LOG_LEVEL)
	var configPath string
	var keyboardFile string
	var port uint
	var appConfig *AppConfig
	var err error

	flag.StringVar(&configPath, "config-path", os.Getenv("SERVER_CONFIG"), "path to the config file (required, attempts to read from 'SERVER_CONFIG' env variable).")
	flag.StringVar(&configPath, "c", os.Getenv("SERVER_CONFIG"), "path to the config file (shorthand).")
	flag.StringVar(&keyboardFile, "keyboard-file", "", fmt.Sprintf("path to the keyboard device. (default '%s')", KEYBOARD_DEFAULT_FILE))
	flag.StringVar(&keyboardFile, "k", "", "path to the keyboard device (shorthand).")
	flag.UintVar(&port, "port", 0, fmt.Sprintf("Port for server to listen on. (default '%n')", server.DEFAULT_PORT))
	flag.UintVar(&port, "p", 0, "Port for server to listen on.")

	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Cwd '%s'\n", cwd)

	flag.Parse()

	logger.Infof("From CLI flag config path '%s'\n", configPath)
	if appConfig, err = LoadConfig(configPath); err != nil {
		logger.Fatal(err)
	}

	logger.LogLevel = log.GetLevel(appConfig.LogLevel)
	logger.Infof("Log level set from config file to: '%s'", logger.LogLevel)

	if keyboardFile == "" {
		keyboardFile = appConfig.Keyboard.File
	}

	logger.Infof("Keyboard file '%s'", keyboardFile)
	f, err := os.OpenFile(keyboardFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	var kf keyboard.File
	kf.File = *f
	kf.StrokeDelay = time.Millisecond * time.Duration(appConfig.Keyboard.StrokeDelayMs)
	defer kf.Close()

	server.New(
		appConfig.Server,
		log.New("Server", appConfig.Server.LogLevel),
		&kf,
	).Run()
}
