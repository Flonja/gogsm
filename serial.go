package gogsm

import (
	"github.com/tarm/serial"
)

func FromSerial(port string) (GSMDevice, error) {
	serialPort, err := serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
	if err != nil {
		return nil, err
	}
	_ = serialPort
	return &DefaultGSMDevice{serialPort}, nil
}
