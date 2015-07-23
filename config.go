package main

import (
	"errors"
	"os"
	"path/filepath"
)

type goingConfig struct {
	Log              logConfig
	PidFile          string
	ProgramConfigDir string
	SocketPath       string
	Umask            int
	Username         string
}

// Search order for config file in this order:
// global variable flagConfigFile
// environment variable named in const EnvVarConfigFile
// $HOME/.going.conf,
// /etc/going.conf
func findGoingConfigFile() (string, error) {
	if isValidFile(*flagConfigFile) {
		return *flagConfigFile, nil
	}

	envVar := os.Getenv(EnvVarConfigFile)
	if isValidFile(envVar) {
		return envVar, nil
	}

	home := os.Getenv("HOME")
	homeConfig := filepath.Join(home, ".going", "going.conf")
	if isValidFile(homeConfig) {
		return homeConfig, nil
	}

	etcConfig := "/etc/going/going.conf"
	if isValidFile(etcConfig) {
		return etcConfig, nil
	}

	return "", errors.New("Could not find config file.")
}

func newGoingConfig(filename string) (*goingConfig, error) {
	c := new(goingConfig)
	err := loadYaml(filename, c)

	configDir := filepath.Dir(filename)

	// set defaults
	if c.PidFile == "" {
		c.PidFile = filepath.Join(configDir, "going.pid")
	}

	if c.ProgramConfigDir == "" {
		c.ProgramConfigDir = filepath.Join(configDir, "programs")
	}

	if c.Log.Dir == "" {
		c.Log.Dir = filepath.Join(configDir, "logs")
	}

	return c, err
}
