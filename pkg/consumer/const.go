package consumer

type Type int

const (
	NSQConsumer Type = iota
	ChannelConsumer
)

func (ct Type) String() string {
	switch ct {
	case NSQConsumer:
		return "NSQConsumer"
	case ChannelConsumer:
		return "ChannelConsumer"
	default:
		return "Unknown"
	}
}
