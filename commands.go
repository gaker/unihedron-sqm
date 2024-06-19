package sqm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Command string

const (
	CmdAveragedReading   Command = "rx"
	CmdUnaveragedReading Command = "ux"
	CmdCalibrationInfo   Command = "cx"
	CmdUnitInfo          Command = "ix"
	CmdReset             Command = "0x19"
)

var ErrInvalidLength = errors.New("invalid length")

func InvalidStartPositionErr(start string) error {
	return fmt.Errorf("invalid reading. '%s' start position not valid", start)
}

// Readable interface is used when reading responses
type Readable interface {
	Parse(in []byte) error
}

// convertNumber is a helper function to convert
// the ASCII data returned from the device into
// an int or a float64
func convertNumber[T int | float64](in []byte) T {
	var zero T

	val := string(in)
	val = strings.TrimSpace(val)

	switch any(zero).(type) {
	case float64:
		if i, err := strconv.ParseFloat(val, 64); err == nil {
			return any(i).(T)
		}
	case int:
		if i, err := strconv.Atoi(val); err == nil {
			return any(i).(T)
		}
	}

	return zero
}

// validStartChar validates the starting character
// read from the unit
func validStartChar(in byte, choices []byte) bool {
	for _, x := range choices {
		if in == x {
			return true
		}
	}

	return false
}

// UnitInfo is the structure for the "rx" or "ux" commands
//
// Example response:
// r,-09.42m,0000005915Hz,0000000000c,0000000.000s, 027.0C
type Reading struct {
	// A reading can be requested as either averaged or unaveraged.
	//
	// request: rx - returns r (averaged)
	// request: ux - returns u (unaveraged)
	Averaged bool `json:"averaged" yaml:"averaged" xml:"averaged"`

	// Reading in mag/arcsec2
	// Leading space for positive value.
	// Leading negative sign (-) for negative value.
	// A reading of 0.00m means that the light at the sensor has reached
	// the upper brightness limit of the unit.
	Reading float64 `json:"reading" yaml:"reading" xml:"reading"`

	// Frequency of sensor in Hz.
	Frequency int `json:"frequency" yaml:"frequency" xml:"frequency"`

	// Period of sensor in counts. Counts occur at a
	// rate of 460.8 kHz (14.7456MHz/32).
	Counts int `json:"counts" yaml:"counts" xml:"counts"`

	// Period of sensor in seconds with millisecond resolution.
	// Determined by dividing the above counts value by 460800.
	Millis float64 `json:"millis" yaml:"millis" xml:"millis"`

	// Temperature measured at light sensor in degrees C.
	// The value is averaged and presented every 4.3 seconds.
	// Leading space for positive value.
	// Leading negative sign (-) for negative value.
	Temp float64 `json:"temp" yaml:"temp" xml:"temp"`
}

// Parse implements the Readable interface
func (r *Reading) Parse(in []byte) error {
	if len(in) != 57 && len(in) != 66 {
		// it can possibly have a serial number on
		// the response as well. For the time being
		// I'm ignoring that.
		return ErrInvalidLength
	}

	if ok := validStartChar(in[0], []byte{'r', 'u'}); !ok {
		return InvalidStartPositionErr(string(in[0]))
	}

	r.Reading = convertNumber[float64](in[2:7])
	r.Frequency = convertNumber[int](in[10:20])
	r.Counts = convertNumber[int](in[23:33])
	r.Millis = convertNumber[float64](in[35:46])
	r.Temp = convertNumber[float64](in[48:54])

	if in[0] == byte('r') {
		r.Averaged = true
	}

	return nil
}

// UnitInfo is the structure for the "ix" command
//
// Example: i,00000002,00000003,00000001,00000413
type UnitInfo struct {
	// Protocol number (8 digits). This will always be the first 8
	// characters (after the “i,” response). This value indicates
	// the revision number of the data protocol to/from the SQM-LE.
	// The protocol version is independent of the feature version.
	Protocol int `json:"protocol" yaml:"protocol" xml:"protocol"`

	// Model number (8 digits). The model value identifies the specific
	// hardware model that the firmware is tailored for.
	Model int `json:"model" yaml:"model" xml:"model"`

	// Feature number (8 digits). The feature value identifies software
	// features. This number is independent of the data protocol.
	Feature int `json:"feature" yaml:"feature" xml:"feature"`

	// Serial number (8 digits). Each unit has its own unique serial number.
	SerialNumber int `json:"serial" yaml:"serial" xml:"serial"`
}

// Parse implements the Readable interface
func (ui *UnitInfo) Parse(in []byte) error {
	if len(in) != 39 {
		return ErrInvalidLength
	}

	if ok := validStartChar(in[0], []byte{'i'}); !ok {
		return InvalidStartPositionErr(string(in[0]))
	}

	ui.Protocol = convertNumber[int](in[2:10])
	ui.Model = convertNumber[int](in[11:19])
	ui.Feature = convertNumber[int](in[20:28])
	ui.SerialNumber = convertNumber[int](in[29:37])

	return nil
}

// CalibrationInfo is for the "cx" command
//
// Example:
//
//	c,00000017.60m,0000000.000s, 039.4C,00000008.71m, 039.4C
type CalibrationInfo struct {
	LightCalibrationOffset    float64 `json:"lightCalibrationOffset" yaml:"lightCalibrationOffset" xml:"lightCalibrationOffset"`
	LightCalibrationTemp      float64 `json:"lightCalibrationTemp" yaml:"lightCalibrationTemp" xml:"lightCalibrationTemp"`
	DarkCalibrationTimePeriod float64 `json:"darkCalibrationTimePeriod" yaml:"darkCalibrationTimePeriod" xml:"darkCalibrationTimePeriod"`
	DarkCalibrationTemp       float64 `json:"darkCalibrationTemp" yaml:"darkCalibrationTemp" xml:"darkCalibrationTemp"`
	Offset                    float64 `json:"offset" yaml:"offset" xml:"offset"`
}

// Parse implements the Readable interface
func (ci *CalibrationInfo) Parse(in []byte) error {
	if len(in) != 58 {
		return ErrInvalidLength
	}

	if ok := validStartChar(in[0], []byte{'c'}); !ok {
		return InvalidStartPositionErr(string(in[0]))
	}

	ci.LightCalibrationOffset = convertNumber[float64](in[2:13])
	ci.LightCalibrationTemp = convertNumber[float64](in[15:26])
	ci.DarkCalibrationTimePeriod = convertNumber[float64](in[28:34])
	ci.DarkCalibrationTemp = convertNumber[float64](in[36:47])
	ci.Offset = convertNumber[float64](in[49:55])

	return nil
}
