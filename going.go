package main

import (
	"ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// main app struct
type going struct {
	config *config
	logger *log.Logger
	// key = yaml base filename without extension
	programs map[string]*program
}

// wish it was named getGoing, but consistency trumps all
func newGoing(configFile string) *going {
	c, err := newConfig(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(RETURN_CONFIG_ERROR)
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
		// TODO log
	}

	for _, file := range programs {
		p = new(program)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			// TODO log this
			continue
		}
		err = yaml.Unmarshal(data, p)
		key := p.Name
		if key == "" {
			key = programName(file)
		}
		g.programs[key] = p
	}
}

func (g *going) scanProgramConfigDir() ([]string, error) {
	programs := new([]string)
	dir := g.config.ProgramConfigDir
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			programs = append(programs, path)
		}
		return nil
	})
	return programs, err
}

func (g *going) loadLogger() {
	var file *os.File
	logFile := filepath.Join(g.config.Log.Path, "going.log")
	if fileExists(logFile) {
		file, _ = os.OpenFile(logFile, O_APPEND, os.ModeAppend)
	} else {
		file, _ = os.Create(logFile)
	}
	g.logger = log.New(file, "", os.LstdFlags)
}
