package parsing

import "fmt"

const (
	MobileEquipmentMessageStorage MessageStorage = iota
	SimMessageStorage
	StatusReportMessageStorage
	AllMessageStorage
)

type MessageStorage uint8

func (c MessageStorage) String() string {
	switch c {
	case MobileEquipmentMessageStorage:
		return "ME"
	case SimMessageStorage:
		return "SM"
	case StatusReportMessageStorage:
		return "SR"
	case AllMessageStorage:
		return "MT"
	}
	panic("shouldn't happen")
}

func MessageStorageFromString(s string) MessageStorage {
	switch EncodedString(s).RemoveQuotes() {
	case "ME":
		return MobileEquipmentMessageStorage
	case "SM":
		return SimMessageStorage
	case "SR":
		return StatusReportMessageStorage
	case "MT":
		return AllMessageStorage
	}
	panic(fmt.Errorf("unknown message storage: %v", s))
}

type MessageStorageUsage struct {
	Current     MessageStorage
	UsedSpace   int
	MaxMessages int
}
