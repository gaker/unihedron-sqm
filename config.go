package sqm

import "errors"

type Tcp struct {
	Host *string `json:"host" yaml:"host" default:"0.0.0.0"`
	Port *string `json:"port" yaml:"port" default:"10001"`
}

type Serial struct {
	Port *string `json:"serialPort" yaml:"serialPort"`
}

type Config struct {
	Tcp    *Tcp    `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	Serial *Serial `json:"serial,omitempty" yaml:"serial,omitempty"`
}

func (c *Config) Validate() error {
	if c.Tcp == nil && c.Serial == nil {
		return errors.New("[SQM] one of http or serial should be set in the config")
	}

	if c.Tcp != nil {
		// we have required fields
		if c.Tcp.Host == nil || c.Tcp.Port == nil {
			return errors.New("tcp.host and tcp.port are required")
		}

		return nil
	}

	if c.Serial.Port == nil {
		return errors.New("serial.port is required")
	}

	return nil
}
