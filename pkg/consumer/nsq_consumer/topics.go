package nsq_consumer

import (
	"sync"
)

var (
	consumerTopics = make(map[string]struct{})
	mux            sync.Mutex
)

// AddConsumerTopic adds a topic to the map if it does not already exist.
func AddConsumerTopic(topic string) {
	mux.Lock()
	defer mux.Unlock()
	consumerTopics[topic] = struct{}{}
}

// RemoveConsumerTopic removes a topic from the map if it exists.
func RemoveConsumerTopic(topic string) {
	mux.Lock()
	defer mux.Unlock()
	delete(consumerTopics, topic)
}

// ExistsConsumerTopic checks if a topic exists in the ConsumerTopics map.
func ExistsConsumerTopic(topic string) bool {
	mux.Lock()
	defer mux.Unlock()
	_, exists := consumerTopics[topic]
	return exists
}

// ConsumerTopics returns a slice of all topics.
func ConsumerTopics() []string {
	mux.Lock()
	defer mux.Unlock()
	var topics []string
	for topic := range consumerTopics {
		topics = append(topics, topic)
	}
	return topics
}
