package pubsub

import "fmt"

var (
	NoSubsMsg       = "No subscribers are connected to the server"
	SubConnectedMsg = "A new subscriber connected to the server"
)

type BrokerConnection interface {
	CreatePub() (chan string, bool)
	CreateSub() chan string
	ReceiveMsg(string)
}

type Broker struct {
	subs []chan string
	pubs []chan string
}

func NewBroker() BrokerConnection {
	return &Broker{
		subs: make([]chan string, 0),
		pubs: make([]chan string, 0),
	}
}

func (b *Broker) CreateSub() chan string {
	fmt.Println("Sub Creation")
	subChan := make(chan string)
	b.subs = append(b.subs, subChan)
	b.notifyPubs()
	fmt.Println("Sub Created")
	return subChan
}

func (b *Broker) notifyPubs() {
	for _, pub := range b.pubs {
		pub <- SubConnectedMsg
	}
}

func (b *Broker) CreatePub() (chan string, bool) {
	fmt.Println("Pub Creation")
	pubChan := make(chan string)
	b.pubs = append(b.pubs, pubChan)
	fmt.Println("Pub Created")
	return pubChan, b.hasSubs()
}

func (b *Broker) hasSubs() bool {
	return len(b.subs) > 0
}

func (b *Broker) ReceiveMsg(msg string) {
	fmt.Println("Message received", msg)
	for _, sub := range b.subs {
		sub <- msg
		fmt.Println("Message Sent")
	}

}
