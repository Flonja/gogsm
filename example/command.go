package main

import (
	"fmt"
	"github.com/flonja/gogsm"
)

func main() {
	gsmDevice, err := gogsm.FromSerial("/dev/ttyUSB2")
	if err != nil {
		return
	}
	resp, err := gsmDevice.SignalQuality()
	if err != nil {
		return
	}
	fmt.Println(resp.DBM.Description()) // Excellent
}
