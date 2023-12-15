package main

import (
	"fmt"
	"sync"

	"github.com/RarePepeCode/quick-broker/pkg/pubsub"
	"github.com/RarePepeCode/quick-broker/pkg/quic"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	broker := pubsub.NewBroker()
	shutdownCh := make(chan string)
	pubErrCh, err := quic.PubConn(broker, shutdownCh)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for pubErr := range pubErrCh {
			fmt.Println(pubErr)
		}
	}()

	subErrCh, err := quic.SubConn(broker)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for subErr := range subErrCh {
			fmt.Println(subErr)
		}
	}()
	wg.Wait()
}
