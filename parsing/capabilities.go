package parsing

import (
	"fmt"
	"strings"
)

const (
	GSMCommandSet CommandSetCapability = iota
	FaxCommandSet
	DataServiceCommandSet
	MobileSpecificCommandSet
)

type CommandSetCapability uint8

func CommandSetFromString(s string) CommandSetCapability {
	switch strings.TrimPrefix(s, "+") {
	case "CGSM":
		return GSMCommandSet
	case "FCLASS":
		return FaxCommandSet
	case "DS":
		return DataServiceCommandSet
	case "MS":
		return MobileSpecificCommandSet
	}
	panic(fmt.Errorf("unknown capability: %v", s))
}
