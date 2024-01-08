package gogsm

import (
	"fmt"
	"github.com/flonja/gogsm/parsing"
)

func (d *DefaultGSMDevice) Check() error {
	_, err := d.ExecuteCommand("AT")
	return err
}

func (d *DefaultGSMDevice) SignalQuality() (parsing.SignalQuality, error) {
	resp, err := d.executeSimpleCommand("+CSQ")
	if err != nil {
		return parsing.SignalQuality{}, err
	}
	sq, err := parsing.SignalQualityString(resp).Parsed()
	if err != nil {
		return parsing.SignalQuality{}, err
	}
	return sq, nil
}

func (d *DefaultGSMDevice) Manufacturer() (string, error) {
	if err := d.testCommand("+CGMI"); err != nil {
		return "", err
	}
	return d.executeSimpleCommand("+CGMI")
}

func (d *DefaultGSMDevice) Model() (string, error) {
	if err := d.testCommand("+CGMM"); err != nil {
		return "", err
	}
	return d.executeSimpleCommand("+CGMM")
}

func (d *DefaultGSMDevice) Revision() (string, error) {
	if err := d.testCommand("+CGMR"); err != nil {
		return "", err
	}
	return mapped(replace("Revision:"), wrap(d.ExecuteCommand("AT+CGMR"))...)
}

func (d *DefaultGSMDevice) SerialNumber() (string, error) {
	if err := d.testCommand("+CGSN"); err != nil {
		return "", err
	}
	return d.executeSimpleCommand("+CGSN")
}

func (d *DefaultGSMDevice) SubscriberId() (string, error) {
	if err := d.testCommand("+CIMI"); err != nil {
		return "", err
	}
	return d.executeSimpleCommand("+CIMI")
}

func (d *DefaultGSMDevice) Capabilities() ([]parsing.CommandSetCapability, error) {
	return mapped(mapArray(parsing.CommandSetFromString), wrap(mapped(
		split(","), wrap(d.executeSimpleCommand("+GCAP"))...))...)
}

func (d *DefaultGSMDevice) ProductIdentification() (string, error) {
	return d.executeSimpleCommand("I")
}

// Utilities:
func (d *DefaultGSMDevice) testCommand(cmd string) error {
	_, err := d.executeSimpleCommand(fmt.Sprintf("%v=?", cmd))
	return err
}

func (d *DefaultGSMDevice) getCommand(cmd string) (string, error) {
	return d.executeSimpleCommand(fmt.Sprintf("%v?", cmd))
}

func (d *DefaultGSMDevice) setCommand(cmd string, value string) (string, error) {
	return d.executeSimpleCommand(fmt.Sprintf("%v=%v", cmd, value))
}

func (d *DefaultGSMDevice) executeSimpleCommand(cmd string) (string, error) {
	return mapped(replace(fmt.Sprintf("%v: ", cmd)), wrap(d.ExecuteCommand(fmt.Sprintf("AT%v", cmd)))...)
}
