package main

import (
	"sync"
	"time"

	"github.com/RarePepeCode/quick-broker/pkg/pubsub"
	"github.com/RarePepeCode/quick-broker/pkg/quic"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	broker := pubsub.NewBroker()
	quic.PubConn(broker)
	quic.SubConn(broker)
	time.Sleep(5 * time.Second)
	wg.Wait()
}
