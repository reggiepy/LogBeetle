package manager

import (
	"context"
	"fmt"
	"sync"
)

type Consumer interface {
	Start() error
	Stop() error

	Topic() string
}

type Manager struct {
	consumers []Consumer
	mux       sync.Mutex
	cnt       int
}

func (m *Manager) Start() {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, c := range m.consumers {
		err := c.Start()
		if err != nil {
			fmt.Println(err)
		}
	}
	return
}

func (m *Manager) Stop() {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, c := range m.consumers {
		err := c.Stop()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (m *Manager) StartWithCtxBackend(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	}
}

func (m *Manager) Add(c Consumer) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.consumers = append(m.consumers, c)
	m.cnt += 1
}

func (m *Manager) Count() int {
	return m.cnt
}

func (m *Manager) Topics() (ret []string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, c := range m.consumers {
		ret = append(ret, c.Topic())
	}
	return ret
}
