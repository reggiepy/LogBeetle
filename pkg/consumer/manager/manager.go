package manager

import (
	"fmt"
	"sync"
)

type Manager struct {
	consumers []LogBeetleConsumer
	mux       sync.Mutex
	cnt       int
}

func (m *Manager) Stop() {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, c := range m.consumers {
		err := c.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (m *Manager) Add(c LogBeetleConsumer) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.consumers = append(m.consumers, c)
	m.cnt += 1
}

func (m *Manager) Count() int {
	return m.cnt
}
