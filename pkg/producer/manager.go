package producer

import (
	"fmt"
	"sync"
)

type Manager struct {
	producers []LogBeetleProducer
	mux       sync.Mutex
}

func (m *Manager) Stop() {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, c := range m.producers {
		err := c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (m *Manager) Add(c LogBeetleProducer) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.producers = append(m.producers, c)
}

func NewLogBeetleProducerManager() *Manager {
	return &Manager{}
}
