package main

import (
	"syscall"
)

type program struct {
	Command     string
	Environment map[string]string
	Log         logConfig
	Name        string
	ProcessName string
	StopSignal  syscall.Signal
	Username    string
	WorkingDir  string
}
