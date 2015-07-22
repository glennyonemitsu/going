package main

import (
	"flag"
	"log"
	"os"
)

const (
	ENV_VAR_CONFIG_FILE string = "GOING_CONFIG_FILE"
)

const (
	RETURN_CONFIG_ERROR int = iota
)

var flagConfigFile *string

func init() {
	flagConfigFile = flag.String("config", "", "Config file")
	flag.Parse()
}

func main() {
	configFile, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
		os.Exit(RETURN_CONFIG_ERROR)
	}
	g := newGoing(configFile)
}
