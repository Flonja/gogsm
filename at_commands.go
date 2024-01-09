package gogsm

import (
	"errors"
	"fmt"
	"github.com/flonja/gogsm/parsing"
	"strconv"
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

func (d *DefaultGSMDevice) CharacterSet() (parsing.CharacterSet, error) {
	resp, err := d.getCommand("+CSCS")
	if err != nil {
		return parsing.GSM7BitCharacterSet, err
	}
	return parsing.CharacterSetFromString(resp), nil
}

func (d *DefaultGSMDevice) SetCharacterSet(set parsing.CharacterSet) error {
	return d.setCommand("+CSCS", fmt.Sprintf(`"%s"`, set))
}

func (d *DefaultGSMDevice) NetworkOperator() (string, error) {
	parts, err := mapped(split(","), wrap(d.getCommand("+COPS"))...)
	if err != nil {
		return "", err
	}
	if len(parts) < 2 {
		return "", errors.New("no operator found")
	}
	return string(parsing.EncodedString(parts[2]).RemoveQuotes()), nil
}

func (d *DefaultGSMDevice) SetPreferredMessageStorage(storage parsing.MessageStorage) error {
	return d.setCommand("+CPMS", fmt.Sprintf(`"%s","%s","%s"`, storage, parsing.SimMessageStorage, parsing.SimMessageStorage))
}

func (d *DefaultGSMDevice) MessageFormat() (parsing.MessageFormat, error) {
	messageFormatRaw, err := d.getCommand("+CMGF")
	if err != nil {
		return parsing.PDUMessageFormat, err
	}
	messageFormat, err := strconv.Atoi(messageFormatRaw)
	if err != nil {
		return parsing.PDUMessageFormat, err
	}
	return parsing.MessageFormat(messageFormat), nil
}

func (d *DefaultGSMDevice) SetMessageFormat(format parsing.MessageFormat) error {
	return d.setCommand("+CMGF", fmt.Sprintf("%d", format))
}

// Utilities:
func (d *DefaultGSMDevice) testCommand(cmd string) error {
	_, err := d.executeCommand(cmd, "=?")
	return err
}

func (d *DefaultGSMDevice) getCommand(cmd string) (string, error) {
	return d.executeCommand(cmd, "?")
}

func (d *DefaultGSMDevice) setCommand(cmd string, value string) error {
	_, err := d.executeCommand(cmd, fmt.Sprintf("=%v", value))
	return err
}

func (d *DefaultGSMDevice) executeSimpleCommand(cmd string) (string, error) {
	return d.executeCommand(cmd, "")
}

func (d *DefaultGSMDevice) executeCommand(cmd string, extras string) (string, error) {
	return mapped(replace(fmt.Sprintf("%v: ", cmd)), wrap(d.ExecuteCommand(fmt.Sprintf("AT%v%v", cmd, extras)))...)
}
