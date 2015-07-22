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
	c, err := getConfig()
	if err != nil {
		log.Fatalf("Could not get config file: %s", err)
		os.Exit(RETURN_CONFIG_ERROR)
	}
}

// main watcher
func going() {
}
