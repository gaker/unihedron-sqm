package sqm

// Monitor interface
//
// There are a types of Sky Quality meters.
//
//  1. SQM-LU USB/Serial - Supported with "Usb"
//  2. SQML-LU-DL Not supported yet
//  3. SQM-LE Ethernet - Supported with "Network"
//  4. SQM-LR RS232 - Not supported yet
type Monitor interface {
	Dial() error
	Send(command Command) error

	// Read a response from the SQM.
	// Note: implementors should close the connection
	// after reading.
	Read(item Readable) error
}

// Creates a new instance of Monitor
// based on values in the configuration
func New(cfg *Config) (Monitor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Add support for other
	// connection types in the future
	return &Network{cfg: cfg}, nil
}
