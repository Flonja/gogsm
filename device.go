package gogsm

import (
	"errors"
	"fmt"
	"github.com/flonja/gogsm/parsing"
	"io"
	"strconv"
	"strings"
	"time"
)

type IncomingSMSMessage struct {
	SmsMessage parsing.SMSMessage
	Storage    parsing.MessageStorage
}

type GSMDevice interface {
	io.Closer

	IncomingSMSMessage() <-chan IncomingSMSMessage

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
	// SetPreferredMessageStorage sets the preferred message storage for incoming SMS/MMS messages.
	SetPreferredMessageStorage(storage parsing.MessageStorage) error
	// MessageStorageAndUsage returns the preferred message storage, used up space and maximum space allowed for incoming SMS/MMS messages.
	MessageStorageAndUsage() (parsing.MessageStorageUsage, error)
	// MessageFormat returns the message format used to encode/decode SMS/MMS messages.
	MessageFormat() (parsing.MessageFormat, error)
	// SetMessageFormat sets the message format used to encode/decode SMS/MMS messages.
	SetMessageFormat(format parsing.MessageFormat) error
	// SMSMessages returns all messages from the provided parsing.MessageFilter in the parsing.MessageStorage.
	SMSMessages(storage parsing.MessageStorage, filter parsing.MessageFilter) ([]parsing.SMSMessage, error)
	// SMSMessage returns a message from the provided index in the parsing.MessageStorage.
	SMSMessage(storage parsing.MessageStorage, index int) (parsing.SMSMessage, error)
	// DeleteAllSMSMessages deletes all messages with the provided parsing.MessageDeleteFilter in the parsing.MessageStorage.
	DeleteAllSMSMessages(storage parsing.MessageStorage, filter parsing.MessageDeleteFilter) error
	// DeleteSMSMessage deletes a message with the provided index in the parsing.MessageStorage.
	DeleteSMSMessage(storage parsing.MessageStorage, index int) error
}

func NewGSMDevice(socket io.ReadWriteCloser) (GSMDevice, error) {
	dev := &DefaultGSMDevice{socket: socket, incomingSMSMessages: make(chan IncomingSMSMessage)}
	if err := dev.Check(); err != nil {
		return nil, err
	}
	go dev.watch()
	return dev, nil
}

type DefaultGSMDevice struct {
	closed bool
	socket io.ReadWriteCloser

	executingCommand    bool
	incomingSMSMessages chan IncomingSMSMessage
}

func (d *DefaultGSMDevice) watch() error {
	for !d.closed {
		time.Sleep(time.Second * 5)
		if d.executingCommand {
			continue
		}

		buf := make([]byte, 256)
		n, err := d.socket.Read(buf)
		if err != nil {
			return err
		}
		out := strings.TrimSpace(string(buf[:n]))
		if argsStr, ok := strings.CutPrefix(out, "+CMTI: "); ok {
			args := strings.Split(argsStr, ",")
			storage := parsing.MessageStorageFromString(args[0])
			index, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			smsMessage, err := d.SMSMessage(storage, index)
			if err != nil {
				return err
			}
			d.incomingSMSMessages <- IncomingSMSMessage{smsMessage, storage}
		} else {
			return fmt.Errorf("unknown command: %v", out)
		}
	}
	return nil
}

func (d *DefaultGSMDevice) Close() error {
	d.closed = true
	return d.socket.Close()
}

func (d *DefaultGSMDevice) WriteString(s string) (n int, err error) {
	if d.closed {
		return 0, errors.New("closed")
	}

	return d.socket.Write([]byte(s))
}

func (d *DefaultGSMDevice) ExecuteCommand(s string) (resp string, err error) {
	d.executingCommand = true
	defer func() {
		d.executingCommand = false
	}()

	if _, err = d.WriteString(fmt.Sprintf("%s\r\n", s)); err != nil {
		return "", err
	}
	out := ""
	for {
		buf := make([]byte, 256)
		n, err := d.socket.Read(buf)
		if err != nil {
			return "", err
		}
		out += string(buf[:n])
		if after, ok := strings.CutSuffix(strings.TrimSpace(out), OK); ok {
			out = after
			break
		}
		if after, ok := strings.CutSuffix(strings.TrimSpace(out), ERROR); ok {
			return "", errors.New(after)
		}
	}
	return strings.TrimSpace(out), err
}

const (
	OK    = "OK"
	ERROR = "ERROR"
)
