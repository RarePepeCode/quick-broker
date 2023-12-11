package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test creation of pubsub broker
func TestBrokerCreation(t *testing.T) {
	broker := NewBroker()
	assert.NotNil(t, broker)
}

// Test cretion of publisher for the server
func TestPubCreation(t *testing.T) {
	broker := NewBroker()
	c, hasSubs := broker.CreatePub()
	assert.NotNil(t, c)
	assert.False(t, hasSubs)
}

// As when creating a pub we receive info if any subs are connected,
// we can check successful sub creation by first creating it and then creating pub.
func TestSubCreation(t *testing.T) {
	broker := NewBroker()
	subChan := broker.CreateSub()
	pubChan, hasSubs := broker.CreatePub()
	assert.NotNil(t, subChan)
	assert.NotNil(t, pubChan)
	assert.True(t, hasSubs)
}

// Creates multiple subs, which should receive a message when broker gets it from the pub.
func TestMessageReceive(t *testing.T) {
	message := "Good message"
	broker := NewBroker()
	c1 := broker.CreateSub()
	c2 := broker.CreateSub()
	go broker.ReceiveMsg(message)
	assert.Equal(t, message, <-c1)
	assert.Equal(t, message, <-c2)
}

// Closes all channels for the server and checks if none of them can receive messages.
func TestClose(t *testing.T) {
	broker := NewBroker()
	c1 := broker.CreateSub()
	c2 := broker.CreateSub()
	c3 := broker.CreateSub()
	c4, _ := broker.CreatePub()

	broker.Close(c2)
	broker.Close(c3)
	broker.Close(c1)
	broker.Close(c4)

	_, ok1 := (<-c1)
	_, ok2 := (<-c2)
	_, ok3 := (<-c3)
	_, ok4 := (<-c4)

	assert.False(t, ok1)
	assert.False(t, ok2)
	assert.False(t, ok3)
	assert.False(t, ok4)

}
