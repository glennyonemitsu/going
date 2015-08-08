package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

func isValidFile(filename string) bool {
	return filename != "" && fileExists(filename)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// programName gets the program name from the filename base without extension
func programName(filename string) string {
	filename = filepath.Base(filename)
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}

func loadYaml(filename string, target interface{}) error {
	if !filepath.IsAbs(filename) {
		filename, _ = filepath.Abs(filename)
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Could not open config file \"%s\": %s", filename, err)
	}
	err = yaml.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("Could not process config file as yaml data: \"%s\", %s", filename, err)
	}
	return nil
}

// exampleConf dumps the exported fields of config structs to directly use in
// yaml files.
func exampleConf(s reflect.Type, prefix string) string {
	output := ""
	for i := 0; i < s.NumField(); i += 1 {
		field := s.Field(i)
		output += fmt.Sprintf(
			"%s %s\n%s %s:\n%s\n",
			prefix,
			field.Tag.Get("cdoc"),
			prefix,
			field.Name,
			prefix,
		)
		if field.Type.Kind() == reflect.Struct {
			output += exampleConf(field.Type, prefix+"    ")
		}
	}
	return output
}
