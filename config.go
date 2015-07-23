package main

import (
	"errors"
	"fmt"
	"ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type goingConfig struct {
	Log              logConfig
	PidFile          string
	PollInterval     int
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
		return *flagConfig
	}

	envVar := os.Getenv(EnvVarConfigFile)
	if isValidFile(envVar) {
		return envVar
	}

	home := os.Getenv("HOME")
	homeConfig := path.Join(home, ".going.conf")
	if isValidFile(homeConfig) {
		return homeConfig
	}

	etcConfig := "/etc/going.conf"
	if isValidFile(etcConfig) {
		return etcConfig
	}

	return "", errors.New("Could not find config file.")
}

func newConfig(filename string) (*goingConfig, error) {
	c := new(config)
	if err != nil {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			err = fmt.Errorf("Could not open config file \"%s\": %s", filename, err)
		}
		err = yaml.Unmarshal(data, c)
		if err != nil {
			err = fmt.Errorf("Could not process config file as yaml data: \"%s\", %s", filename, err)
		}
	}
	return c, err
}
