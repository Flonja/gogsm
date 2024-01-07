package parsing

import (
	"encoding/hex"
	"github.com/warthog618/sms/encoding/ucs2"
	"strconv"
	"strings"
	"unicode/utf16"
)

type EncodedString string

func (s EncodedString) FromUCS2HexString() (EncodedString, error) {
	hexString, err := hex.DecodeString(string(s))
	if err != nil {
		return "", err
	}
	decoded, err := ucs2.Decode(hexString)
	if err != nil {
		return "", err
	}
	return EncodedString(decoded), nil
}

func (s EncodedString) FromAsciiString() (EncodedString, error) {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] < 2 && i+2 < len(s) {
			uInt16, err := strconv.ParseUint(string(s[i])+string(s[i+1])+string(s[i+2]), 0, 16)
			if err != nil {
				return "", err
			}
			for _, r := range utf16.Decode([]uint16{uint16(uInt16)}) {
				b.WriteRune(r)
			}
			i += 3
		} else if i+1 < len(s) {
			uInt16, err := strconv.ParseUint(string(s[i])+string(s[i+1]), 0, 16)
			if err != nil {
				return "", err
			}
			for _, r := range utf16.Decode([]uint16{uint16(uInt16)}) {
				b.WriteRune(r)
			}
			i += 2
		} else {
			i += 1
		}
	}
	return EncodedString(b.String()), nil
}

func (s EncodedString) ToUCS2HexString() (EncodedString, error) {
	encoded := ucs2.Encode([]rune(s))
	hexString := hex.EncodeToString(encoded)
	return EncodedString(hexString), nil
}

func (s EncodedString) ToAsciiString() (s1 EncodedString) {
	for _, i := range utf16.Encode([]rune(s)) {
		s1 += EncodedString(strconv.Itoa(int(i)))
	}
	return
}
