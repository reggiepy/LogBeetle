package producer

type LogBeetleProducer interface {
	Close() error
	Publish(topic string, body []byte) error
}
