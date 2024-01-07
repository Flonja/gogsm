package parsing

import (
	"strconv"
	"strings"
)

type SignalQuality struct {
	DBM DecibelMilliwatts
	BER BitErrorRate
}

type SignalQualityString string

func (s SignalQualityString) Parsed() (sq SignalQuality, err error) {
	signal := strings.Split(string(s), ",")
	signalStrength, err := strconv.Atoi(signal[0])
	if err != nil {
		return
	}
	ber, err := strconv.Atoi(signal[0])
	if err != nil {
		return
	}
	sq.BER = BitErrorRate(ber)

	if signalStrength != 99 {
		sq.DBM = DecibelMilliwatts(-113 + signalStrength*2)
	}
	return
}

type DecibelMilliwatts int

func (dbm DecibelMilliwatts) Description() string {
	if dbm >= -70 {
		return "Excellent"
	}
	if dbm >= -85 {
		return "Good"
	}
	if dbm >= -100 {
		return "Fair"
	}
	if dbm > -110 {
		return "Poor"
	}
	return "No Signal"
}

type BitErrorRate int

func (ber BitErrorRate) Description() string {
	switch ber {
	case 0:
		return "less than 0.2%"
	case 1:
		return "0.2% to 0.4%"
	case 2:
		return "0.4% to 0.8%"
	case 3:
		return "0.8% to 1.6%"
	case 4:
		return "1.6% to 3.2%"
	case 5:
		return "3.2% to 6.4%"
	case 6:
		return "6.4% to 12.8%"
	case 7:
		return "more than 12.8%"
	}
	return "N/A"
}
