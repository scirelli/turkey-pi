package main

import (
	"flag"
	"os"

	"github.com/scirelli/turkey-pi/internal/app/server"
	"github.com/scirelli/turkey-pi/pkg/keyboard"
	"github.com/scirelli/turkey-pi/pkg/log"
)

func main() {
	var logger = log.New("Main", log.DEFAULT_LOG_LEVEL)
	var configPath string
	var keyboardFile string
	var appConfig *AppConfig
	var err error

	flag.StringVar(&configPath, "config-path", os.Getenv("SERVER_CONFIG"), "path to the config file.")
	flag.StringVar(&configPath, "c", os.Getenv("SERVER_CONFIG"), "path to the config file (shorthand).")
	flag.StringVar(&keyboardFile, "keyboard-file", "/dev/hidg0", "path to the keyboard device.")
	flag.StringVar(&keyboardFile, "k", "/dev/hidg0", "path to the keyboard device (shorthand).")

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

	logger.Infof("Keyboard file '%s'", keyboardFile)
	f, err := os.OpenFile(keyboardFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	var kf keyboard.File
	kf.File = *f
	defer kf.Close()

	if _, err := kf.WriteString("steve\n"); err != nil {
		logger.Fatal(err)
	}

	server.New(
		appConfig.Server,
		log.New("Server", appConfig.Server.LogLevel),
		&kf,
	).Run()
}
