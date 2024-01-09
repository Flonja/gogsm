package parsing

const (
	NationalPhoneNumberType      PhoneNumberType = 129
	InternationalPhoneNumberType PhoneNumberType = 145
	TextPhoneNumberType          PhoneNumberType = 208
	Text2PhoneNumberType         PhoneNumberType = 209 // I know, it's stupid...
)

type PhoneNumberType uint8

func (p PhoneNumberType) String() string {
	switch p {
	case NationalPhoneNumberType:
		return "National"
	case InternationalPhoneNumberType:
		return "International"
	case TextPhoneNumberType:
		fallthrough
	case Text2PhoneNumberType:
		return "Text-based"
	}
	panic("shouldn't happen")
}

func (p PhoneNumberType) IsText() bool {
	return p == TextPhoneNumberType || p == Text2PhoneNumberType
}
