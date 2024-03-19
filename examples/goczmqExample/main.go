package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4/draft"
)

func main() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5555")
	subscriber.SetSubscribe("") // 订阅所有消息

	for {
		msg, _ := subscriber.Recv(0)
		fmt.Println("Received:", msg)
	}
}
