package parsing

const (
	PDUMessageFormat MessageFormat = iota
	TextMessageFormat
)

type MessageFormat uint8
