package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
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
	c, err := newGoingConfig(configFile)
	if err != nil {
		log.Print(err)
		os.Exit(ReturnConfigError)
	}
	g := new(going)
	g.programs = make(map[string]*program)
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
		p.state = StateNil
		p.config, err = newProgramConfig(file)
		if err != nil {
			g.logger.Printf("Could not read program config file %s. Skipping.", file)
			continue
		}
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
			g.logger.Printf("Loading program conf file %s.", path)
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
		file, _ = os.OpenFile(logFile, os.O_WRONLY, 0644)
	} else {
		file, _ = os.Create(logFile)
	}
	g.logger = log.New(file, "", log.LstdFlags)
}

func (g *going) runPrograms() {
	for _, program := range g.programs {
		go g.runProgram(program)
	}
}

func (g *going) runProgram(p *program) {
	name := p.config.Name
	g.logger.Printf("Initializing program %s", name)
	err := p.init()
	if err != nil {
		g.logger.Printf("Could not initializing program %s. Error: %s", name, err)
	}
	g.logger.Printf("Running program %s", name)
	// proc := p.cmd.Process
	//state := p.cmd.ProcessState
	err = p.cmd.Run()
	if err != nil {
		g.logger.Printf("Could not run program %s. Error: %s", name, err)
	}
	for {
		time.Sleep(500 * time.Millisecond)
	}
}

func (g *going) listen() {
	for {
		time.Sleep(500 * time.Millisecond)
	}
}
