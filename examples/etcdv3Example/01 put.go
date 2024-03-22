package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Println(1, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Put(ctx, "sample_key", time.Now().String())
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
	// use the response
	defer cli.Close()
}
