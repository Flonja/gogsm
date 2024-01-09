package parsing

import (
	"fmt"
)

const (
	ReceivedUnreadMessageFilter MessageFilter = iota
	ReceivedReadMessageFilter
	StoredUnsentMessageFilter
	StoredSentMessageFilter
	AllMessageFilter
)

type MessageFilter uint8

func (c MessageFilter) String() string {
	switch c {
	case ReceivedUnreadMessageFilter:
		return "REC UNREAD"
	case ReceivedReadMessageFilter:
		return "REC READ"
	case StoredUnsentMessageFilter:
		return "STO UNSENT"
	case StoredSentMessageFilter:
		return "STO SENT"
	case AllMessageFilter:
		return "ALL"
	}
	panic("shouldn't happen")
}

func MessageFilterFromString(s string) MessageFilter {
	switch EncodedString(s).RemoveQuotes() {
	case "REC UNREAD":
		return ReceivedUnreadMessageFilter
	case "REC READ":
		return ReceivedReadMessageFilter
	case "STO UNSENT":
		return StoredUnsentMessageFilter
	case "STO SENT":
		return StoredSentMessageFilter
	case "ALL":
		return AllMessageFilter
	}
	panic(fmt.Errorf("unknown message filter: %v", s))
}
