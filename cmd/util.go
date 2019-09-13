package cmd

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	cmdBlockWrite    = 0x07
	cmdBlockRead     = 0x06
	cmdPageErase     = 0x08
	cmdDeviceErase   = 0x03
	cmdGetVersion    = 0x01
	cmdGetDerivative = 0x02
	cmdDirectRead    = 0x09
	cmdDirectWrite   = 0x0A
	cmdIndirectRead  = 0x0B
	cmdIndirectWrite = 0x0C
)

func erasePage(number uint8) error {
	logrus.WithFields(logrus.Fields{
		"Page": number,
	}).Debug("erasing page")

	// FPDAT is hard-coded here for EFM8SB1
	fpdat := uint8(0xB4)

	// Write the FPDAT address
	if err := prog.WriteAR(fpdat); err != nil {
		return err
	}

	if err := prog.WriteCommand(cmdPageErase); err != nil {
		return err
	}

	res, err := prog.ReadData()
	if err != nil {
		return err
	}

	if res != 0x0D {
		return fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	// Write page number
	if err := prog.WriteCommand(number); err != nil {
		return err
	}

	res, err = prog.ReadData()
	if err != nil {
		return err
	}

	if res != 0x0D {
		return fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	// Write 0x00
	if err := prog.WriteCommand(0x00); err != nil {
		return err
	}

	res, err = prog.ReadData()
	if err != nil {
		return err
	}

	if res != 0x0D {
		return fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	return nil
}

func writeBlock(addr uint16, data []uint8) error {
	logrus.WithFields(logrus.Fields{
		"Address": addr,
		"Data":    data,
	}).Debug("writing block")

	// FPDAT is hard-coded here for EFM8SB1
	fpdat := uint8(0xB4)

	// Write the FPDAT address
	if err := prog.WriteAR(fpdat); err != nil {
		return err
	}

	if err := prog.WriteCommand(cmdBlockWrite); err != nil {
		return err
	}

	res, err := prog.ReadData()
	if err != nil {
		return err
	}

	if res != 0x0D {
		return fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	// Write high byte of address
	if err := prog.WriteCommand(uint8(addr >> 8)); err != nil {
		return err
	}

	// Write low byte of address
	if err := prog.WriteCommand(uint8(addr)); err != nil {
		return err
	}

	// Write the length
	if err := prog.WriteCommand(uint8(len(data))); err != nil {
		return err
	}

	// Write the data
	for i := 0; i < len(data); i++ {
		// Write the byte
		if err := prog.WriteCommand(data[i]); err != nil {
			return err
		}
	}

	res, err = prog.ReadData()
	if err != nil {
		return err
	}

	if res != 0x0D {
		return fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	return nil
}

func readBlock(addr uint16, length uint8) ([]uint8, error) {
	logrus.WithFields(logrus.Fields{
		"Address": addr,
		"Length":  length,
	}).Debug("reading block")

	// FPDAT is hard-coded here for EFM8SB1
	fpdat := uint8(0xB4)

	// Write the FPDAT address
	if err := prog.WriteAR(fpdat); err != nil {
		return nil, err
	}

	if err := prog.WriteCommand(cmdBlockRead); err != nil {
		return nil, err
	}

	res, err := prog.ReadData()
	if err != nil {
		return nil, err
	}

	if res != 0x0D {
		return nil, fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	// Write high byte of address
	if err := prog.WriteCommand(uint8(addr >> 8)); err != nil {
		return nil, err
	}

	// Write low byte of address
	if err := prog.WriteCommand(uint8(addr)); err != nil {
		return nil, err
	}

	// Write the length
	if err := prog.WriteCommand(length); err != nil {
		return nil, err
	}

	res, err = prog.ReadData()
	if err != nil {
		return nil, err
	}

	if res != 0x0D {
		return nil, fmt.Errorf("unexpected command result: expected: 0x0D, actual 0x%x", res)
	}

	buf := make([]byte, length)
	for i := 0; i < int(length); i++ {
		b, err := prog.ReadData()
		if err != nil {
			return nil, err
		}

		buf[i] = b
	}

	return buf, nil
}

func initChip() {
	if !prog.Check() {
		logrus.Fatal("programmer does not answer")
	}

	if err := prog.Halt(); err != nil {
		logrus.WithError(err).Fatal("failed to halt chip")
	}

	<-time.After(time.Millisecond * 50)

	if err := prog.WriteAR(0x00); err != nil {
		logrus.WithError(err).Fatal("failed to select device id register")
	}

	id, err := prog.ReadDR()
	if err != nil {
		logrus.WithError(err).Fatal("failed to read chip id")
	}

	logrus.WithFields(logrus.Fields{
		"Chip ID": id,
	}).Info("Found chip")

	if id != 52 {
		// For other devices, implement the init procedure (AN127 p. 27) and the flashing procedure.
		logrus.Fatal("Only EFM8LB1 is currently supported!")
	}

	// init the EFM8SB1 (see AN127)
	if err := prog.WriteSFR(0xB6, 0x40); err != nil {
		logrus.WithError(err).Fatal("chip init failed - flash timimg")
	}

	if err := prog.WriteSFR(0xFF, 0x80); err != nil {
		logrus.WithError(err).Fatal("chip init failed - vdd mon 1")
	}
	//delay 5us
	<-time.After(time.Microsecond * 5)
	if err := prog.WriteSFR(0xEF, 0x02); err != nil {
		logrus.WithError(err).Fatal("chip init failed - vdd mon 2")
	}

	if err := prog.WriteSFR(0xA9, 0x00); err != nil {
		logrus.WithError(err).Fatal("chip init failed - osc")
	}

	<-time.After(time.Millisecond * 50)
}
