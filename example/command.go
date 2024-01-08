package main

import (
	"fmt"
	"github.com/flonja/gogsm"
)

func main() {
	gsmDevice := must(gogsm.FromSerial("/dev/ttyUSB2"))
	resp := must(gsmDevice.SignalQuality())
	fmt.Println(resp.DBM.Description()) // Excellent
	fmt.Printf("%#v\n", must(gsmDevice.Model()))
	fmt.Printf("%#v\n", must(gsmDevice.Manufacturer()))
	fmt.Printf("%#v\n", must(gsmDevice.Revision()))
	fmt.Printf("%#v\n", must(gsmDevice.ProductIdentification()))
	fmt.Printf("%#v\n", must(gsmDevice.SubscriberId()))
	fmt.Printf("%#v\n", must(gsmDevice.Capabilities()))
	fmt.Printf("%v\n", must(gsmDevice.CharacterSet()))
	fmt.Printf("%v\n", must(gsmDevice.NetworkOperator()))
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
