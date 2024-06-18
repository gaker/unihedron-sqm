package sqm

// r,-09.42m,0000005915Hz,0000000000c,0000000.000s, 027.0C
type Reading struct{}

type Monitor interface {
	Read(cfg *Config) error
}

type Usb struct{}

func (u *Usb) Read(cfg *Config) error {
	return nil
}

type Network struct{}

func (n *Network) Read(cfg *Config) error {
	return nil
}

// Creates a new instance of Monitor
// based on values in the configuration
func New(cfg *Config) (Monitor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if cfg.Tcp != nil {
		return &Network{}, nil
	}

	return &Usb{}, nil
}
