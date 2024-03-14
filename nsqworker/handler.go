package nsqworker

import (
	"github.com/nsqio/go-nsq"
)

type MessageHandler struct {
	Handler func(message []byte) error
}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	err := h.Handler(m.Body)
	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}
