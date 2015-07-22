package main

import (
	"os"
)

func isValidFile(filename string) bool {
	return filename != "" && fileExists(filename)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsExist(err)
}
