package cmd

import (
	"io/ioutil"
	"math"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var firmware string
var erase bool

var flashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Flash a new firmware to the device",
	Run: func(cmd *cobra.Command, args []string) {
		initChip()

		fw, err := ioutil.ReadFile(firmware)
		if err != nil {
			logrus.WithError(err).Fatal("failed to read firmware file")
		}

		// pad to 128 bytes
		for i := 0; i < 128-(len(fw)%128); i++ {
			fw = append(fw, 0xFF)
		}

		if erase {
			// page size is hard coded here
			count := uint8(math.Ceil(float64(len(fw)) / 512))
			for i := uint8(0); i < count; i++ {
				logrus.Debugf("Erasing page: %d", i)
				if err := erasePage(i); err != nil {
					logrus.WithError(err).Fatal("failed to erase page")
				}
			}
		}

		bar := pb.StartNew(len(fw))

		for i := 0; i < len(fw); i += 128 {
			end := i + 128
			if end > len(fw) {
				end = len(fw)
			}

			if err := writeBlock(uint16(i), fw[i:end]); err != nil {
				logrus.WithError(err).Fatal("failed to write block")
			}

			bar.Add(128)
		}

		prog.Reset()
	},
}

func init() {
	flashCmd.Flags().StringVar(&firmware, "fw", "fw.bin", "Firmware file")
	flashCmd.Flags().BoolVar(&erase, "auto-erase", true, "Auto erase")
	rootCmd.AddCommand(flashCmd)
}
