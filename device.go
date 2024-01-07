package gogsm

import "io"

type GSMDevice interface {
	io.ReadWriter
	io.StringWriter
}

func NewGSMDevice(socket io.ReadWriter) *DefaultGSMDevice {
	return &DefaultGSMDevice{socket: socket}
}

type DefaultGSMDevice struct {
	socket io.ReadWriter
}

func (d DefaultGSMDevice) Read(p []byte) (n int, err error) {
	return d.socket.Read(p)
}

func (d DefaultGSMDevice) Write(p []byte) (n int, err error) {
	return d.socket.Write(p)
}

func (d DefaultGSMDevice) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// TODO: implement proper command execution
