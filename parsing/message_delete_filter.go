package parsing

const (
	// ReadMessageDeleteFilter Delete all read messages from storage,
	// leaving unread messages and stored mobile originated messages (whether sent or not) untouched
	ReadMessageDeleteFilter MessageDeleteFilter = iota + 1
	// ReadAndSentMessageDeleteFilter Delete all read messages from storage and sent mobile originated messages,
	// leaving unread messages and unsent mobile originated messages untouched
	ReadAndSentMessageDeleteFilter
	// ReadSentAndUnsentMessageFilter Delete all read messages from storage, sent and unsent mobile originated
	// messages, leaving unread messages untouched
	ReadSentAndUnsentMessageFilter
	// AllMessageDeleteFilter Self explanatory: deletes all messages from storage.
	AllMessageDeleteFilter
)

type MessageDeleteFilter uint8
