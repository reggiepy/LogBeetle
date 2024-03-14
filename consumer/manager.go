package consumer

import (
	"sync"
)

var (
	consumers []*LogConsumer
	mu        sync.Mutex
)

func StopConsumers() {
	mu.Lock()
	defer mu.Unlock()
	for _, c := range consumers {
		c.Close()
	}
}

func AddConsumer(c *LogConsumer) {
	mu.Lock()
	defer mu.Unlock()
	consumers = append(consumers, c)
}
