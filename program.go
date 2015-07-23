package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
)

const (
	StateNil = iota
	StateInitialized
	StateRunning
)

type program struct {
	config *programConfig
	// internal program logger to capture all program output
	logger *log.Logger
	// actual user struct data based on Username lookup
	user  *user.User
	state int
	cmd   *exec.Cmd
}

type programConfig struct {
	Command     string
	Environment map[string]string
	Log         logConfig
	Name        string
	ProcessName string
	StopSignal  syscall.Signal
	Username    string
	Dir         string
}

// init sets up other data and structs and performs checks to make sure
// everything is available to properly run. The logger param is sent by the
// parent going struct.
func (p *program) init() error {
	var err error

	p.user, err = user.Lookup(p.config.Username)
	if err != nil {
		return fmt.Errorf(
			"Could not lookup username \"%s\" for program \"%s\".",
			p.config.Username,
			p.config.Name,
		)
	}

	err = p.loadLogger()
	if err != nil {
		return fmt.Errorf(
			"Could not load logger for program \"%s\".",
			p.config.Name,
		)
	}

	p.cmd = exec.Command(p.config.Command)
	p.cmd.Dir = p.config.Dir
	for key, value := range p.config.Environment {
		p.cmd.Env = append(p.cmd.Env, key+"="+value)
	}

	p.state = StateInitialized

	return err
}

func (p *program) run() {
	p.cmd.Run()

}

func (p *program) loadLogger() error {
	var file *os.File
	var err error
	logFile := filepath.Join(p.config.Log.Dir, p.config.Name+".log")
	if fileExists(logFile) {
		file, err = os.OpenFile(logFile, os.O_APPEND, os.ModeAppend)
	} else {
		file, err = os.Create(logFile)
	}
	if err != nil {
		return err
	}
	p.logger = log.New(file, "", log.LstdFlags)
	return nil
}
