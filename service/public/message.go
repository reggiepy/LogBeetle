package public

import (
	"fmt"
	"time"

	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/arrayUtils"
	"github.com/reggiepy/LogBeetle/pkg/consumer"
)

type Message struct{}

func (m *Message) SendMessage(projectName string, message string) error {
	if !arrayUtils.InArray(projectName, consumer.Topics) {
		return fmt.Errorf("topic【%s】 is not allowed\n", projectName)
	}
	start := time.Now()
	// 向 NSQ 发送消息
	err := global.LbNsqProducer.Publish(projectName, []byte(message))
	if err != nil {
		return fmt.Errorf("topic【%s】send message faild: %v\n", projectName, err)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	global.LbLogger.Info(fmt.Sprintf("topic【%s】消息写入时间：%s\n", projectName, elapsed))
	return nil
}
