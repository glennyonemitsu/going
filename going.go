package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// main app struct
type going struct {
	config *goingConfig
	logger *log.Logger
	// key = yaml base filename without extension
	programs map[string]*program
}

// wish it was named getGoing, but consistency trumps all
func newGoing(configFile string) *going {
	c, err := newConfig(configFile)
	if err != nil {
		log.Print(err)
		os.Exit(ReturnConfigError)
	}
	g := new(going)
	g.config = c
	g.loadLogger()
	g.getPrograms()
	return g
}

func (g *going) getPrograms() {
	var p *program
	var programs []string
	var err error

	programs, err = g.scanProgramConfigDir()
	if err != nil {
		g.logger.Printf("Could not scan program config files in %s.", g.config.ProgramConfigDir)
		os.Exit(ReturnProgramScanError)
	}

	for _, file := range programs {
		p = new(program)
		p.config = new(programConfig)
		p.state = StateNil
		data, err := ioutil.ReadFile(file)
		if err != nil {
			g.logger.Printf("Could not read program config file %s. Skipping.", file)
			continue
		}
		err = yaml.Unmarshal(data, p.config)
		key := p.config.Name
		if key == "" {
			key = programName(file)
		}
		g.programs[key] = p
	}
}

func (g *going) scanProgramConfigDir() ([]string, error) {
	programs := []string{}
	dir := g.config.ProgramConfigDir
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".conf" {
			programs = append(programs, path)
		}
		return nil
	})
	return programs, err
}

func (g *going) loadLogger() {
	var file *os.File
	logFile := filepath.Join(g.config.Log.Dir, "going.log")
	if fileExists(logFile) {
		file, _ = os.OpenFile(logFile, os.O_APPEND, os.ModeAppend)
	} else {
		file, _ = os.Create(logFile)
	}
	g.logger = log.New(file, "", log.LstdFlags)
}

func (g *going) runPrograms() {
	for name, program := range g.programs {
		g.logger.Printf("Initializing program %s", name)
		err := program.init()
		if err != nil {
			g.logger.Printf("Could not initializing program %s. Error: %s", name, err)
			continue
		}
		g.logger.Printf("Running program %s", name)
		go program.run()
	}
}
