package sqm

import (
	"errors"
	"fmt"
)

// TCP settings
type Tcp struct {
	Host *string `json:"host" yaml:"host" default:"0.0.0.0"`
	Port *string `json:"port" yaml:"port" default:"10001"`
}

func (t *Tcp) Addr() string {
	if t.Host == nil && t.Port == nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", *t.Host, *t.Port)
}

// Config
type Config struct {
	Tcp *Tcp `json:"tcp,omitempty" yaml:"tcp,omitempty"`
}

// Validate runs simple validation against the config struct
func (c *Config) Validate() error {
	if c.Tcp != nil {
		// we have required fields
		if c.Tcp.Host == nil || c.Tcp.Port == nil {
			return errors.New("tcp.host and tcp.port are required")
		}

		return nil
	}

	return errors.New("tcp configuration required")
}
