package cmd

import "github.com/spf13/cobra"
import "github.com/Sirupsen/logrus"

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Used for debugging",
	Run: func(cmd *cobra.Command, args []string) {
		if !prog.Check() {
			logrus.Fatal("programmer does not answer")
		}

		if err := prog.Reset(); err != nil {
			logrus.WithError(err).Fatal("failed to reset chip")
		}

		if err := prog.WriteAR(0x01); err != nil {
			logrus.WithError(err).Fatal("failed to write the address register")
		}

		addr, err := prog.ReadAR()
		if err != nil {
			logrus.WithError(err).Fatal("failed to read the address register")
		}

		if addr != 0x01 {
			logrus.Fatalf("address register read wrong: expected: 0x01, actual: 0x%x", addr)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
