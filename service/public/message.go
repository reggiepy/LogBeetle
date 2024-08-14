package public

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/consumer"
	"github.com/reggiepy/LogBeetle/pkg/producer"
	"github.com/reggiepy/LogBeetle/util/array_utils"
	"go.uber.org/zap"
	"time"
)

type Message struct{}

func (m *Message) SendMessage(projectName string, message string) error {
	if !array_utils.InArray(projectName, consumer.Topics) {
		return fmt.Errorf("topic【%s】 is not allowed\n", projectName)
	}
	start := time.Now()
	// 向 NSQ 发送消息
	err := producer.Instance.Publish(projectName, []byte(message))
	if err != nil {
		return fmt.Errorf("topic【%s】send message faild: %v\n", projectName, err)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	zap.L().Info(fmt.Sprintf("topic【%s】消息写入时间：%s\n", projectName, elapsed))
	return nil
}
