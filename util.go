package main

import (
	"os"
	"path"
)

func isValidFile(filename string) bool {
	return filename != "" && fileExists(filename)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsExist(err)
}

// programName gets the program name from the filename base without extension
func programName(filename string) string {
	filename = path.Base(filename)
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
