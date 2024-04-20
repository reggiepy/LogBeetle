package nsqconsumer

import (
	"sync"

	"github.com/nsqio/go-nsq"
)

var (
	nsqs []*nsq.Consumer
	mu   sync.Mutex
)

func StopConsumers() {
	mu.Lock()
	defer mu.Unlock()
	for _, c := range nsqs {
		c.Stop()
	}
}

func AddConsumer(c *nsq.Consumer) {
	mu.Lock()
	defer mu.Unlock()
	nsqs = append(nsqs, c)
}
