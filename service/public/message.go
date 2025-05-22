package public

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
)

type ServiceMessage struct{}

func (m *ServiceMessage) SendMessage(projectName string, message string) error {
	if !global.LBConsumerManager.ExistName(projectName) {
		return fmt.Errorf("topic【%s】 is not allowed\n", projectName)
	}
	// 向 NSQ 发送消息
	err := global.LbNsqProducer.Publish(projectName, []byte(message))
	if err != nil {
		return fmt.Errorf("topic【%s】send message faild: %v\n", projectName, err)
	}
	return nil
}
