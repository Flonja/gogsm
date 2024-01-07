package gogsm

import (
	"github.com/flonja/gogsm/parsing"
	"strings"
)

func (d *DefaultGSMDevice) Check() error {
	_, err := d.ExecuteCommand("AT")
	return err
}

func (d *DefaultGSMDevice) SignalQuality() (parsing.SignalQuality, error) {
	resp, err := d.ExecuteCommand("AT+CSQ")
	if err != nil {
		return parsing.SignalQuality{}, err
	}
	sq, err := parsing.SignalQualityString(strings.TrimPrefix(resp, "+CSQ: ")).Parsed()
	if err != nil {
		return parsing.SignalQuality{}, err
	}
	return sq, nil
}
