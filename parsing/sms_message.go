package parsing

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type SMSMessagesString string

func (s SMSMessagesString) Parsed(prefix string) (messages []SMSMessage, err error) {
	list := strings.Split(string(s), "\r\n")
	for i := 0; i < len(list); i += 2 {
		parts := strings.Split(strings.TrimPrefix(
			strings.ReplaceAll(list[i], "\n", ""), prefix+": "), ",")
		msg := SMSMessage{}

		msg.Index, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		msg.Status = MessageFilterFromString(parts[1])

		sender, err := EncodedString(parts[2]).RemoveQuotes().FromUCS2HexString()
		if err != nil {
			return nil, err
		}
		numberType, err := strconv.Atoi(parts[6])
		if err != nil {
			return nil, err
		}
		msg.PhoneNumberType = PhoneNumberType(numberType)
		if msg.PhoneNumberType.IsText() {
			sender, err = sender.FromAsciiString()
			if err != nil {
				return nil, err
			}
		}
		msg.Sender = string(sender)

		date := strings.TrimPrefix(parts[4], `"`)
		timeAndZone := strings.TrimSuffix(parts[5], `"`)
		msg.Time, err = parseTime(date + " " + timeAndZone)
		if err != nil {
			return nil, err
		}

		textLength, err := strconv.Atoi(parts[7])
		if err != nil {
			return nil, err
		}
		messageText := strings.ReplaceAll(list[i+1], "\n", "")
		hexString, err := hex.DecodeString(messageText)
		if err != nil {
			return nil, err
		}
		if len(hexString)/textLength == 2 {
			msgText, err := EncodedString(messageText).FromUCS2HexString()
			if err != nil {
				return nil, err
			}
			msg.Text = string(msgText)
		} else {
			msg.Text = string(hexString)
		}
		messages = append(messages, msg)
	}
	return
}

func parseTime(s string) (time.Time, error) {
	timeParts := strings.Split(s, "+")
	timezoneNumber, err := strconv.ParseFloat(timeParts[1], 64)
	if err != nil {
		return time.Time{}, err
	}
	timezone := timezoneNumber / 4.0
	timezonePrefix := "-"
	if timezone > 0 {
		timezonePrefix = "+"
	}
	return time.Parse("06/01/02 15:04:05-07", timeParts[0]+fmt.Sprintf("%v%02.0f", timezonePrefix, math.Abs(timezone)))
}

type SMSMessage struct {
	Index           int
	Status          MessageFilter
	PhoneNumberType PhoneNumberType
	Sender          string
	Time            time.Time
	Text            string
}
