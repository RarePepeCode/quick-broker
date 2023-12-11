package pubsub

import (
	"fmt"
	"slices"
)

var (
	NoSubsMsg       = "No subscribers are connected to the server"
	SubConnectedMsg = "A new subscriber connected to the server"
)

// Interface that allow pub and sub peer communincate with each other
type BrokerConnection interface {

	// Broker establishes a new connection with the publisher peer.
	// Returns channel through which messages will be sent to publisher client.
	// When creating a new publisher server will return if there's any connected subscribers at the moment
	CreatePub() (chan string, bool)

	// Broker establishes a new connection with the subscriber peer.
	// Returns channel through which messages will be sent to subscriber client.
	CreateSub() chan string

	// Function through which publisher can send message to all exisiting subscribers.
	ReceiveMsg(string)

	// Closes the connection to given channel, as well deletes it from the pool of pubs/subs.
	// Messages can no longer pass throguh this channel.
	Close(chan string)
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

func (b *Broker) Close(c chan string) {
	var index int
	if slices.Contains(b.subs, c) {
		index = slices.Index(b.subs, c)
		b.subs = append(b.subs[:index], b.subs[index+1:]...)
	} else if slices.Contains(b.pubs, c) {
		index = slices.Index(b.pubs, c)
		b.subs = append(b.pubs[:index], b.pubs[index+1:]...)
	}
	close(c)
}
