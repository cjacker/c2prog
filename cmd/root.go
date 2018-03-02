package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/donothingloop/c2prog/programmer"
	"github.com/spf13/cobra"
)

var debug bool
var prog programmer.Programmer

var rootCmd = &cobra.Command{
	Use:               "c2prog",
	Short:             "c2prog - Programmer for the C2 protocol of 8051 based Silicon Labs chips",
	PersistentPreRun:  start,
	PersistentPostRun: stop,
}

func start(cmd *cobra.Command, args []string) {
	p, err := programmer.NewC2Prog()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open the programmer port")
		return
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	}

	prog = p
}

func stop(cmd *cobra.Command, args []string) {
	if err := prog.Close(); err != nil {
		logrus.WithError(err).Warn("Failed to close programmer port")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
}

// Execute the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("failed to execute")
	}
}
