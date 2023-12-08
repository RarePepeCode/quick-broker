package main

import (

	//"github.com/quic-go/quic-go"

	"fmt"
	"sync"

	"github.com/RarePepeCode/quick-broker/pkg/pubsub"
	"github.com/RarePepeCode/quick-broker/pkg/quic"
	//"github.com/quic-go/quic-go"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("Program started")
	broker := pubsub.NewBroker()
	fmt.Println("Broker created")
	quic.PubConn(broker)
	fmt.Println("Pub Connection started")
	quic.SubConn(broker)
	fmt.Println("Sub Connection started")
	quic.ClientMain()
	wg.Wait()
}
