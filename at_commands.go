package gogsm

import (
	"errors"
	"fmt"
	"github.com/flonja/gogsm/parsing"
	"strconv"
	"strings"
)

func (d *DefaultGSMDevice) IncomingSMSMessage() <-chan IncomingSMSMessage {
	return d.incomingSMSMessages
}

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

func (d *DefaultGSMDevice) MessageStorageAndUsage() (parsing.MessageStorageUsage, error) {
	usage, err := d.getCommand("+CPMS")
	if err != nil {
		return parsing.MessageStorageUsage{}, err
	}
	params := strings.Split(usage, ",")
	used, err := strconv.Atoi(params[1])
	if err != nil {
		return parsing.MessageStorageUsage{}, err
	}
	maxMessages, err := strconv.Atoi(params[1])
	if err != nil {
		return parsing.MessageStorageUsage{}, err
	}
	return parsing.MessageStorageUsage{Current: parsing.MessageStorageFromString(params[0]), UsedSpace: used, MaxMessages: maxMessages}, nil
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

func (d *DefaultGSMDevice) SMSMessages(storage parsing.MessageStorage, filter parsing.MessageFilter) ([]parsing.SMSMessage, error) {
	if err := d.setupSMSMessages(storage); err != nil {
		return nil, err
	}
	messages, err := d.ExecuteCommand(fmt.Sprintf(`AT+CMGL="%v"`, filter))
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, nil
	}
	smsMessages, err := parsing.SMSMessagesString(messages).Parsed("+CMGL")
	if err != nil {
		return nil, err
	}
	return smsMessages, nil
}

func (d *DefaultGSMDevice) SMSMessage(storage parsing.MessageStorage, index int) (parsing.SMSMessage, error) {
	if err := d.setupSMSMessages(storage); err != nil {
		return parsing.SMSMessage{}, err
	}
	messages, err := d.ExecuteCommand(fmt.Sprintf(`AT+CMGR=%v`, index))
	if err != nil {
		return parsing.SMSMessage{}, err
	}
	if len(messages) == 0 {
		return parsing.SMSMessage{}, nil
	}
	messages = strings.Replace(messages, "+CMGR: ", fmt.Sprintf("+CMGR: %v,", index), 1) // Small hack, I know
	smsMessages, err := parsing.SMSMessagesString(messages).Parsed("+CMGR")
	if err != nil {
		return parsing.SMSMessage{}, err
	}
	return smsMessages[0], nil
}

func (d *DefaultGSMDevice) setupSMSMessages(storage parsing.MessageStorage) error {
	if err := d.SetMessageFormat(parsing.TextMessageFormat); err != nil {
		return err
	}
	if err := d.SetCharacterSet(parsing.UCS2CharacterSet); err != nil {
		return err
	}
	if err := d.SetPreferredMessageStorage(storage); err != nil {
		return err
	}
	return d.setCommand("+CSDH", fmt.Sprintf("%d", 1))
}

func (d *DefaultGSMDevice) DeleteAllSMSMessages(storage parsing.MessageStorage, filter parsing.MessageDeleteFilter) error {
	if err := d.SetPreferredMessageStorage(storage); err != nil {
		return err
	}
	return d.setCommand("+CMGD", fmt.Sprintf("%d,%d", 0, filter))
}

func (d *DefaultGSMDevice) DeleteSMSMessage(storage parsing.MessageStorage, index int) error {
	if err := d.SetPreferredMessageStorage(storage); err != nil {
		return err
	}
	return d.setCommand("+CMGD", fmt.Sprintf("%d", index))
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
