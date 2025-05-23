package nsq_consumer

import "go.uber.org/zap"

type zapNSQLogger struct {
	logger *zap.Logger
}

func (z *zapNSQLogger) Output(callDepth int, s string) error {
	z.logger.Sugar().Info(s) // 可换成 Debug, Warn, Error 等
	return nil
}
