package cmd

import (
	"io/ioutil"
	dbg "runtime/debug"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var file string

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read the content of the flash of the device",
	Run: func(cmd *cobra.Command, args []string) {
		dbg.SetGCPercent(0)

		// initialize the chip
		initChip()

		bar := pb.StartNew(8192)
		gbuf := []byte{}
		for i := 0; i < 8192; i += 128 {
			buf, err := readBlock(uint16(i), 128)
			if err != nil {
				logrus.WithError(err).Fatal("failed to read block")
			}

			gbuf = append(gbuf, buf...)
			bar.Add(128)
		}

		ioutil.WriteFile(file, gbuf, 0666)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVar(&file, "fw", "fw.bin", "Firmware file")
}
