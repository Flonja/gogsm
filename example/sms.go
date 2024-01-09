package main

import (
	"fmt"
	"github.com/flonja/gogsm"
	"github.com/flonja/gogsm/parsing"
	"time"
)

func main() {
	gsmDevice := mustComplete(gogsm.FromSerial("/dev/ttyUSB2"))
	// From +11223344556 (International): Sent at 2024-01-09 18:47:07
	// Hi
	for _, message := range mustComplete(gsmDevice.SMSMessages(parsing.SimMessageStorage, parsing.AllMessageFilter)) {
		fmt.Printf("From %v (%v): Sent at %v\n", message.Sender, message.PhoneNumberType, message.Time.Format(time.DateTime))
		fmt.Println(message.Text)
		fmt.Println()
	}
}

func mustComplete[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
