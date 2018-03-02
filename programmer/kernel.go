package programmer

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

// Magic constants for the c2 interface.
const (
	c2progOpPing         = 1
	c2progOpPong         = 2
	c2progOpReset        = 3
	c2progOpWriteAR      = 4
	c2progOpReadAR       = 5
	c2progOpWriteDR      = 6
	c2progOpReadDR       = 7
	c2progOpPoll         = 8
	c2progOpWriteSFR     = 9
	c2progOpReadSFR      = 10
	c2progOpWriteCmd     = 11
	c2progOpReadResponse = 12
	c2progOpReadData     = 13
	c2progOpHalt         = 14
	c2progOpDevID        = 15
)

const path = "/dev/c2prog"

// C2Prog is a C2 programmer implemented on the Raspberry PI GPIO pins.
type C2Prog struct {
	dev *os.File
}

// NewC2Prog creates a new instance of the Raspberry PI C2 programmer.
func NewC2Prog() (Programmer, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &C2Prog{
		dev: f,
	}, nil
}

type kernelMsg struct {
	op    uint8
	data  uint8
	data2 uint8
}

func (k *kernelMsg) bytes() []byte {
	return []byte{k.op, k.data, k.data2}
}

func kernelMsgFromBuf(buf []byte) *kernelMsg {
	return &kernelMsg{
		op:    buf[0],
		data:  buf[1],
		data2: buf[2],
	}
}

func (r *C2Prog) cmd(m *kernelMsg) (*kernelMsg, error) {
	cbuf := m.bytes()
	_, err := r.dev.Write(cbuf)
	if err != nil {
		return nil, err
	}

	// read reply
	buf := make([]byte, 3)
	_, err = r.dev.Read(buf)
	if err != nil {
		return nil, err
	}

	return kernelMsgFromBuf(buf), nil
}

// Check if the programmer responds to commands.
func (r *C2Prog) Check() bool {
	res, err := r.cmd(&kernelMsg{
		op:    c2progOpPing,
		data:  0,
		data2: 0,
	})

	if err != nil {
		logrus.WithError(err).Warn("failed to communicate with c2prog kernel module")
		return false
	}

	return (res.op == c2progOpPong)
}

// Reset the target.
func (r *C2Prog) Reset() error {
	logrus.Debug("C2Prog: resetting target")
	_, err := r.cmd(&kernelMsg{
		op:    c2progOpReset,
		data:  0,
		data2: 0,
	})

	return err
}

// WriteAR writes the address register.
func (r *C2Prog) WriteAR(addr uint8) error {
	logrus.WithField("Address", addr).Debug("C2Prog: writing address")
	_, err := r.cmd(&kernelMsg{
		op:    c2progOpWriteAR,
		data:  addr,
		data2: 0,
	})

	return err
}

// ReadAR reads the address register.
func (r *C2Prog) ReadAR() (uint8, error) {
	res, err := r.cmd(&kernelMsg{
		op:    c2progOpReadAR,
		data:  0,
		data2: 0,
	})

	if err != nil {
		return 0, err
	}

	return res.data, nil
}

// WriteDR writes the data register.
func (r *C2Prog) WriteDR(data uint8) error {
	logrus.WithField("Data", data).Debug("C2Prog: writing data")
	_, err := r.cmd(&kernelMsg{
		op:    c2progOpWriteDR,
		data:  data,
		data2: 0,
	})

	return err
}

// ReadDR reads the data register.
func (r *C2Prog) ReadDR() (uint8, error) {
	res, err := r.cmd(&kernelMsg{
		op:    c2progOpReadDR,
		data:  0,
		data2: 0,
	})

	if err != nil {
		return 0, err
	}

	return res.data, nil
}

// Poll the given flag.
func (r *C2Prog) Poll(flag uint8, res uint8) error {
	_, err := r.cmd(&kernelMsg{
		op:    c2progOpPoll,
		data:  flag,
		data2: res,
	})

	return err
}

// Halt the device.
func (r *C2Prog) Halt() error {
	_, err := r.cmd(&kernelMsg{
		op:    c2progOpHalt,
		data:  0,
		data2: 0,
	})

	return err
}

// WriteSFR writes an SFR.
func (r *C2Prog) WriteSFR(addr uint8, data uint8) error {
	logrus.WithFields(logrus.Fields{
		"Address": addr,
		"Data":    data,
	}).Debug("C2Prog: writing SFR")

	_, err := r.cmd(&kernelMsg{
		op:    c2progOpWriteSFR,
		data:  addr,
		data2: data,
	})

	return err
}

// ReadSFR reads an SFR.
func (r *C2Prog) ReadSFR(addr uint8) (uint8, error) {
	logrus.WithFields(logrus.Fields{
		"Address": addr,
	}).Debug("C2Prog: reading SFR")

	res, err := r.cmd(&kernelMsg{
		op:    c2progOpReadSFR,
		data:  addr,
		data2: 0,
	})

	if err != nil {
		return 0, err
	}

	return res.data, nil
}

// WriteCommand writes a command to the PI.
func (r *C2Prog) WriteCommand(cmd uint8) error {
	logrus.WithField("Command", cmd).Debug("C2Prog: writing command")
	_, err := r.cmd(&kernelMsg{
		op:    c2progOpWriteCmd,
		data:  cmd,
		data2: 0,
	})

	return err
}

// ReadResponse reads the response to a command.
func (r *C2Prog) ReadResponse() (uint8, error) {
	res, err := r.cmd(&kernelMsg{
		op:    c2progOpReadResponse,
		data:  0,
		data2: 0,
	})

	if err != nil {
		return 0, err
	}

	return res.data, nil
}

// ReadData reads the command data byte.
func (r *C2Prog) ReadData() (uint8, error) {
	res, err := r.cmd(&kernelMsg{
		op:    c2progOpReadData,
		data:  0,
		data2: 0,
	})

	if err != nil {
		return 0, err
	}

	return res.data, nil
}

// Close the programmer instance.
func (r *C2Prog) Close() error {
	return r.dev.Close()
}

// WriteDirect allows writes to SFRs on devices that have SFR paging.
func (r *C2Prog) WriteDirect(addr uint8, data uint8) error {
	// FPDAT is hard-coded here for EFM8SB1
	fpdat := uint8(0xB4)

	// AddressWrite(FPDAT)
	if err := r.WriteAR(fpdat); err != nil {
		return err
	}

	// WriteCommand(0x0A) - Direct write
	if err := r.WriteCommand(0x0A); err != nil {
		return err
	}

	// ReadData() - 0x0D indicates success, all other return values are errors
	res, err := r.ReadData()
	if err != nil {
		return err
	}

	if res != 0x0D {
		return fmt.Errorf("read error: %d", res)
	}

	// WriteCommand(addr)
	if err := r.WriteCommand(addr); err != nil {
		return err
	}

	// WriteCommand(0x01)
	if err := r.WriteCommand(0x01); err != nil {
		return err
	}

	// WriteCommand(data)
	return r.WriteCommand(data)
}

// ReadDirect allows reads from SFRs on devices that have SFR paging.
func (r *C2Prog) ReadDirect(addr uint8) (uint8, error) {
	// FPDAT is hard-coded here for EFM8SB1
	fpdat := uint8(0xB4)

	// AddressWrite(FPDAT)
	if err := r.WriteAR(fpdat); err != nil {
		return 0, err
	}

	// WriteCommand(0x09) - Direct read
	if err := r.WriteCommand(0x09); err != nil {
		return 0, err
	}

	// ReadData() - 0x0D indicates success, all other return values are errors
	res, err := r.ReadData()
	if err != nil {
		return 0, err
	}

	if res != 0x0D {
		return 0, fmt.Errorf("read error: %d", res)
	}

	// WriteCommand(addr)
	if err := r.WriteCommand(addr); err != nil {
		return 0, err
	}

	// WriteCommand(0x01)
	if err := r.WriteCommand(0x01); err != nil {
		return 0, err
	}

	res, err = r.ReadData()
	if err != nil {
		return 0, err
	}

	return res, nil
}
