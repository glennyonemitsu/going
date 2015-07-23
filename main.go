package main

import (
	"flag"
	"log"
	"os"
)

const (
	EnvVarConfigFile string = "GOING_CONFIG_FILE"
)

const (
	_ = iota
	ReturnConfigError
	ReturnProgramScanError
)

type logConfig struct {
	Interval string
	Limit    int
	Dir      string
}

var flagConfigFile *string

func init() {
	flagConfigFile = flag.String("config", "", "Config file")
	flag.Parse()
}

func main() {
	configFile, err := findGoingConfigFile()
	if err != nil {
		log.Print(err)
		os.Exit(ReturnConfigError)
	}
	g := newGoing(configFile)
	g.runPrograms()
	g.listen()
}
