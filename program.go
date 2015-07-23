package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"syscall"
)

type program struct {
	config *programConfig
	// internal program logger to capture all program output
	logger *log.Logger
	// actual user struct data based on Username lookup
	user *user.User
}

type programConfig struct {
	Command     string
	Environment map[string]string
	Log         logConfig
	Name        string
	ProcessName string
	StopSignal  syscall.Signal
	Username    string
	WorkingDir  string
}

// init sets up other data and structs and performs checks to make sure
// everything is available to properly run. The logger param is sent by the
// parent going struct.
func (p *program) init() error {
	var err error

	p.user, err = user.Lookup(p.Username)
	if err != nil {
		return fmt.Errorf(
			"Could not lookup username \"%s\" for program \"%s\".",
			p.Username,
			p.Name,
		)
	}

	return err
}

func (p *program) run() {

}

func (p *program) loadLogger() {
	var file *os.File
	logFile := filepath.Join(p.Log.Path, p.Name+".log")
	if fileExists(logFile) {
		file, _ = os.OpenFile(logFile, O_APPEND, os.ModeAppend)
	} else {
		file, _ = os.Create(logFile)
	}
	p.logger = log.New(file, "", os.LstdFlags)
}
