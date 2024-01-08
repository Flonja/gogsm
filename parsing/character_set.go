package parsing

import (
	"fmt"
)

const (
	GSM7BitCharacterSet CharacterSet = iota
	UCS2CharacterSet
	InternationalReferenceCharacterSet
)

type CharacterSet uint8

func (c CharacterSet) String() string {
	switch c {
	case GSM7BitCharacterSet:
		return "GSM"
	case UCS2CharacterSet:
		return "UCS2"
	case InternationalReferenceCharacterSet:
		return "IRA"
	}
	panic("shouldn't happen")
}

func CharacterSetFromString(s string) CharacterSet {
	switch EncodedString(s).RemoveQuotes() {
	case "GSM":
		return GSM7BitCharacterSet
	case "UCS2":
		return UCS2CharacterSet
	case "IRA":
		return InternationalReferenceCharacterSet
	}
	panic(fmt.Errorf("unknown character set: %v", s))
}
