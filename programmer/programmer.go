package programmer

// Programmer provides the basic functions of the C2 protocol.
type Programmer interface {
	// Check checks if the programmer is reachable.
	Check() bool

	// Reset the device.
	Reset() error

	// WriteAR writes the address register.
	WriteAR(addr uint8) error

	// ReadAR reads the address register.
	ReadAR() (uint8, error)

	// WriteDR writes the data register.
	WriteDR(data uint8) error

	// ReadDR reads the data register.
	ReadDR() (uint8, error)

	// Poll the given flag.
	Poll(flag uint8, res uint8) error

	// Halt the device.
	Halt() error

	// WriteSFR writes an SFR.
	WriteSFR(addr uint8, data uint8) error

	// ReadSFR reads an SFR.
	ReadSFR(addr uint8) (uint8, error)

	// WriteCommand writes a command to the PI.
	WriteCommand(cmd uint8) error

	// ReadResponse reads the response to a command.
	ReadResponse() (uint8, error)

	// ReadData reads the command data byte.
	ReadData() (uint8, error)

	// Close the programmer instance.
	Close() error

	// WriteDirect allows writes to SFRs on devices that have SFR paging.
	WriteDirect(addr uint8, data uint8) error

	// ReadDirect allows reads from SFRs on devices that have SFR paging.
	ReadDirect(addr uint8) (uint8, error)
}
