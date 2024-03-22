package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {
	const (
		dialTimeout    = 5 * time.Second
		requestTimeout = 10 * time.Second
	)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Println(1, err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// count keys about to be deleted
	gresp, err := cli.Get(ctx, "sample_key", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gresp)
	// delete the keys
	dresp, err := cli.Delete(ctx, "sample_key", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all keys:", int64(len(gresp.Kvs)) == dresp.Deleted)
}
