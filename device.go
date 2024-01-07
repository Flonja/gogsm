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

	// Check sends a basic "AT" command to see if the device is operational
	Check() error
	// SignalQuality sends a basic "AT+CSQ" command to
	// report rssi (https://en.wikipedia.org/wiki/RSSI) and ber (https://en.wikipedia.org/wiki/Bit_error_rate)
	SignalQuality() (parsing.SignalQuality, error)
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
