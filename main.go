package main

import (
	"log"
	"os"

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

	cmdRoot.AddCommand(cmdRun)
	cmdRoot.Execute()
}
