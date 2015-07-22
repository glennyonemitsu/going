package main

import (
	"errors"
	"fmt"
	"ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

type config struct {
	Log              logConfig
	PidFile          string
	PollInterval     int
	ProgramConfigDir string
	SocketPath       string
	Umask            int
	Username         string
}

type logConfig struct {
	Interval string
	Limit    int
	Dir      string
}

// Search order for config file in this order:
// global variable flagConfigFile
// environment variable named in const ENV_VAR_CONFIG_FILE
// $HOME/.going.conf,
// /etc/going.conf
func findConfigFile() (string, error) {
	if isValidFile(*flagConfigFile) {
		return *flagConfig
	}

	envVar := os.Getenv(ENV_VAR_CONFIG_FILE)
	if isValidFile(envVar) {
		return envVar
	}

	home := os.Getenv("HOME")
	homeConfig := path.Join(home, ".going.yaml")
	if isValidFile(homeConfig) {
		return homeConfig
	}

	etcConfig := "/etc/going.yaml"
	if isValidFile(etcConfig) {
		return etcConfig
	}

	return "", errors.New("Could not find config file.")
}

func newConfig(filename string) (*config, error) {
	c := new(config)
	if err != nil {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			err = fmt.Errorf("Could not open yaml config file \"%s\": %s", filename, err)
		}
		err = yaml.Unmarshal(data, c)
		if err != nil {
			err = fmt.Errorf("Could not process yaml config file \"%s\": %s", filename, err)
		}
	}
	return c, err
}
