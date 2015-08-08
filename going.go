package main

import (
	"errors"
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
	programs map[string]*Program
}

type goingConfig struct {
	Log              logConfig
	PidFile          string
	ProgramConfigDir string
	SocketPath       string
	Umask            int
	Username         string
}

// wish it was named getGoing, but consistency trumps all
func newGoing(configFile string) *going {
	c, err := newGoingConfig(configFile)
	if err != nil {
		log.Print(err)
		os.Exit(ReturnConfigError)
	}
	g := new(going)
	g.programs = make(map[string]*Program)
	g.config = c
	g.loadLogger()
	g.getPrograms()
	return g
}

func (g *going) getPrograms() {
	var p *Program
	var programs []string
	var err error

	programs, err = g.scanProgramConfigDir()
	if err != nil {
		g.logger.Printf("Could not scan program config files in %s.", g.config.ProgramConfigDir)
		os.Exit(ReturnProgramScanError)
	}

	for _, file := range programs {
		p = new(Program)
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

func (g *going) runProgram(p *Program) {
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

// Search order for config file in this order:
// global variable flagConfigFile, aka command line -config flag
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
