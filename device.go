package gogsm

import (
	"errors"
	"fmt"
	"github.com/flonja/gogsm/parsing"
	"io"
	"strings"
)

type GSMDevice interface {
	io.ReadWriter
	io.StringWriter

	// Check sends a basic "AT" command to see if the device is operational.
	Check() error
	// SignalQuality returns the current signal quality by values:
	// RSSI (https://en.wikipedia.org/wiki/RSSI) and BER (https://en.wikipedia.org/wiki/Bit_error_rate).
	SignalQuality() (parsing.SignalQuality, error)
	// Model returns the model of the GSM module.
	Model() (string, error)
	// Manufacturer returns the manufacturer of the GSM module.
	Manufacturer() (string, error)
	// Revision returns the revision of the GSM module.
	Revision() (string, error)
	// SerialNumber returns the serial number of the GSM module.
	SerialNumber() (string, error)
	// SubscriberId returns the IMSI (International Mobile Subscriber Identity) of the SIM inserted into the GSM module.
	SubscriberId() (string, error)
	// ProductIdentification sends a basic "ATI" command for product identification information.
	ProductIdentification() (string, error)
	// Capabilities returns a list of the capabilities this GSM module may have.
	Capabilities() ([]parsing.CommandSetCapability, error)
	// CharacterSet returns the current character set selected.
	CharacterSet() (parsing.CharacterSet, error)
	// SetCharacterSet sets the current character set to the provided one.
	SetCharacterSet(set parsing.CharacterSet) error
	// NetworkOperator returns the current network operator providing service to the GSM module.
	NetworkOperator() (string, error)
}

func NewGSMDevice(socket io.ReadWriter) (GSMDevice, error) {
	dev := &DefaultGSMDevice{socket: socket}
	if err := dev.Check(); err != nil {
		return nil, err
	}
	return dev, nil
}

type DefaultGSMDevice struct {
	socket io.ReadWriter
}

func (d *DefaultGSMDevice) Read(p []byte) (n int, err error) {
	return d.socket.Read(p)
}

func (d *DefaultGSMDevice) Write(p []byte) (n int, err error) {
	return d.socket.Write(p)
}

func (d *DefaultGSMDevice) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

func (d *DefaultGSMDevice) ExecuteCommand(s string) (resp string, err error) {
	if _, err = d.WriteString(fmt.Sprintf("%s\r\n", s)); err != nil {
		return "", err
	}
	out := ""
	for {
		buf := make([]byte, 256)
		n, err := d.Read(buf)
		if err != nil {
			return "", err
		}
		out += strings.TrimSpace(string(buf[:n]))
		if after, ok := strings.CutSuffix(out, OK); ok {
			out = after
			break
		}
		if after, ok := strings.CutSuffix(out, ERROR); ok {
			return "", errors.New(after)
		}
		out += "\n"
	}
	return strings.TrimSpace(out), err
}

const (
	OK    = "OK"
	ERROR = "ERROR"
)
