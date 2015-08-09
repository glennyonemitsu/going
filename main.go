package main

import (
	"log"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

const (
	EnvVarConfigFile string = "GOING_CONFIG_FILE"
)

const (
	_ = iota
	ReturnConfigError
	ReturnProgramScanError
)

type logConfig struct {
	Interval string
	Limit    int
	Dir      string
}

var flagConfigFile *string

func main() {
	cmdRoot := &cobra.Command{Use: "going"}

	cmdRun := &cobra.Command{
		Use:   "run",
		Short: "Run the server to watch processes",
		Run: func(cmd *cobra.Command, args []string) {
			configFile, err := findGoingConfigFile()
			if err != nil {
				log.Print(err)
				os.Exit(ReturnConfigError)
			}

			g := newGoing(configFile)
			g.runPrograms()
			g.listen()
		},
	}
	flagConfigFile = cmdRun.Flags().StringP("config", "c", "", "Config file to use")

	cmdExampleConfig := &cobra.Command{
		Use:   "example-config",
		Short: "Dump example config files to standard out",
	}
	cmdExampleConfigGoing := &cobra.Command{
		Use:   "going",
		Short: "Dump example config file for main going server",
		Run: func(cmd *cobra.Command, args []string) {
			println(ExampleConfigGoing)
		},
	}
	cmdExampleConfigProgram := &cobra.Command{
		Use:   "program",
		Short: "Dump example config file for programs running under the going server",
		Run: func(cmd *cobra.Command, args []string) {
			println(ExampleConfigProgram)
		},
	}

	cmdRoot.AddCommand(cmdRun)
	cmdRoot.AddCommand(cmdExampleConfig)
	cmdExampleConfig.AddCommand(cmdExampleConfigGoing, cmdExampleConfigProgram)
	cmdRoot.Execute()
}
