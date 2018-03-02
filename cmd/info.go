package cmd

import "github.com/spf13/cobra"
import "github.com/Sirupsen/logrus"

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about the target chip",
	Run: func(cmd *cobra.Command, args []string) {
		if !prog.Check() {
			logrus.Fatal("programmer does not answer")
		}

		if err := prog.Reset(); err != nil {
			logrus.WithError(err).Fatal("failed to reset chip")
		}

		id, err := prog.ReadDR()
		if err != nil {
			logrus.WithError(err).Fatal("failed to read chip id")
		}

		rev, err := prog.ReadSFR(0x01)
		if err != nil {
			logrus.WithError(err).Fatal("failed to read the revision id")
		}

		sid, err := prog.ReadSFR(0x00)
		if err != nil {
			logrus.WithError(err).Fatal("failed to read the chip id")
		}

		logrus.Infof("Chip ID: %d, Rev ID: %d, Second Chip ID: %d", id, rev, sid)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
