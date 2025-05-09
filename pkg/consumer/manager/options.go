package manager

import "go.uber.org/zap"

type Options func(m *Manager) error

func WithLogger(logger *zap.Logger) Options {
	return func(m *Manager) error {
		m.logger = logger
		return nil
	}
}
