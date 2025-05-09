package manager

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/consumer"
	"go.uber.org/zap"
	"sync"
)

// Manager 管理消费者的结构体
type Manager struct {
	consumers map[string]consumer.Consumer
	mux       sync.Mutex
	logger    *zap.Logger
}

// NewManager 创建一个新的 Manager 实例
func NewManager(opts ...Options) (*Manager, error) {
	m := &Manager{
		consumers: make(map[string]consumer.Consumer), // 初始化消费者映射
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	if m.logger == nil {
		return nil, fmt.Errorf("no logger provided")
	}
	return m, nil
}

// Start 启动所有消费者
func (m *Manager) Start() {
	m.mux.Lock()
	defer m.mux.Unlock()

	// 遍历所有消费者并启动
	for name, c := range m.consumers {
		err := c.Start()
		if err != nil {
			m.logger.Info(fmt.Sprintf("%s consumer start error: %v\n", name, err))
		} else {
			m.logger.Info(fmt.Sprintf("%s consumer start success\n", name))
		}
	}
}

// Stop 停止所有消费者
func (m *Manager) Stop() {
	m.mux.Lock()
	defer m.mux.Unlock()

	// 遍历所有消费者并停止
	for name, c := range m.consumers {
		err := c.Stop()
		if err != nil {
			m.logger.Info(fmt.Sprintf("%s consumer stop error: %v\n", name, err))
		} else {
			m.logger.Info(fmt.Sprintf("%s consumer stop success\n", name))
		}
	}
}

// Add 向 Manager 中添加一个消费者
func (m *Manager) Add(c consumer.Consumer) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	// 检查消费者是否已经存在
	if _, ok := m.consumers[c.GetName()]; ok {
		// 如果已经存在，返回错误
		return fmt.Errorf("%s consumer already exists", c.GetName())
	}

	// 添加新的消费者，并增加计数
	m.consumers[c.GetName()] = c

	return nil
}

// ConsumerCount 返回当前消费者的数量
func (m *Manager) ConsumerCount() int {
	m.mux.Lock()
	defer m.mux.Unlock()

	return len(m.consumers)
}

// Delete 从 Manager 中删除指定名称的消费者
func (m *Manager) Delete(name string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, exists := m.consumers[name]; exists {
		delete(m.consumers, name)
	}
}

// ExistName 检查消费者是否存在
func (m *Manager) ExistName(name string) bool {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, exists := m.consumers[name]
	return exists
}

// GetTypes 返回所有消费者的类型，并去重
func (m *Manager) GetTypes() []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	uniqueTypes := make(map[string]struct{})
	// 初始化容量为0，动态增长
	result := make([]string, 0)

	// 遍历所有消费者并收集其类型
	for _, c := range m.consumers {
		consumerType := c.GetType()
		// 去重逻辑
		if _, exists := uniqueTypes[consumerType]; !exists {
			uniqueTypes[consumerType] = struct{}{}
			result = append(result, consumerType)
		}
	}

	return result
}

// GetNamesByType 根据消费者类型返回所有消费者的名称
func (m *Manager) GetNamesByType(consumerType string) []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	names := make([]string, 0)
	// 遍历所有消费者并根据类型收集其名称
	for name, c := range m.consumers {
		if c.GetType() == consumerType {
			names = append(names, name)
		}
	}

	return names
}
