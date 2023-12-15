package pubsub

import (
	"slices"
	"sync"
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
	lock sync.Mutex
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
	b.lock.Lock()
	defer b.lock.Unlock()
	subChan := make(chan string)
	b.subs = append(b.subs, subChan)
	b.notifyPubs()
	return subChan
}

func (b *Broker) notifyPubs() {
	for _, pub := range b.pubs {
		pub <- SubConnectedMsg
	}
}

func (b *Broker) CreatePub() (chan string, bool) {
	b.lock.Lock()
	defer b.lock.Unlock()
	pubChan := make(chan string)
	b.pubs = append(b.pubs, pubChan)
	return pubChan, b.hasSubs()
}

func (b *Broker) hasSubs() bool {
	return len(b.subs) > 0
}

func (b *Broker) ReceiveMsg(msg string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for _, sub := range b.subs {
		sub <- msg
	}
}

func (b *Broker) Close(c chan string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	var index int
	if slices.Contains(b.subs, c) {
		index = slices.Index(b.subs, c)
		b.subs = append(b.subs[:index], b.subs[index+1:]...)
	} else if slices.Contains(b.pubs, c) {
		index = slices.Index(b.pubs, c)
		b.pubs = append(b.pubs[:index], b.pubs[index+1:]...)
	}
	close(c)
}
