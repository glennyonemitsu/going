package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

const (
	StateNil = iota
	StateInitialized
	StateRunning
)

type Program struct {
	config *ProgramConfig
	state  int
	cmd    *exec.Cmd
}

// ProgramConfig holds all the parameters for going to run and maintain the
// state of the supervised process. These variables are all configurable via
// the program's yaml format config file.
type ProgramConfig struct {
	Command     string
	Name        string
	Environment map[string]string
	StopSignal  syscall.Signal
	Dir         string
}

// init sets up other data and structs and performs checks to make sure
// everything is available to properly run. The logger param is sent by the
// parent going struct.
func (p *Program) init() error {
	var err error

	p.cmd = exec.Command(p.config.Command)
	p.cmd.Dir = p.config.Dir
	for key, value := range p.config.Environment {
		p.cmd.Env = append(p.cmd.Env, key+"="+value)
	}

	p.state = StateInitialized

	return err
}

func (p *Program) run() {
	p.cmd.Run()

}

func newProgramConfig(filename string) (*ProgramConfig, error) {
	c := new(ProgramConfig)
	err := loadYaml(filename, c)

	if err != nil {
		err = fmt.Errorf("Could not load program config file \"%s\": %s", filename, err)
	}
	return c, err
}
