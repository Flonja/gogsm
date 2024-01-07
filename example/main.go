package main

import (
	"fmt"
	"github.com/flonja/gogsm/parsing"
)

func main() {
	signalQuality, err := parsing.SignalQualityString("17,99").Parsed()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", signalQuality)
	fmt.Printf("dBm: %v, ber: %v\n", signalQuality.DBM.Description(), signalQuality.BER.Description())
	encodedString := parsing.EncodedString("00480069")
	asciiString, err := encodedString.FromUCS2HexString()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", asciiString)
}
