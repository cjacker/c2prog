package cmd

import (
	"github.com/spf13/cobra"
)

var eraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "Erase the first page of the device",
	Run: func(cmd *cobra.Command, args []string) {
		// initialize the chip
		initChip()
		erasePage(0)
	},
}

func init() {
	rootCmd.AddCommand(eraseCmd)
}
