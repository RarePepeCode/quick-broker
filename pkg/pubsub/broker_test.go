package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrokerCreation(t *testing.T) {
	broker := NewBroker()
	assert.NotNil(t, broker)
}

func TestPubCreation(t *testing.T) {
	broker := NewBroker()
	c, hasSubs := broker.CreatePub()
	assert.NotNil(t, c)
	assert.False(t, hasSubs)
}

func TestSubCreation(t *testing.T) {
	broker := NewBroker()
	subChan := broker.CreateSub()
	pubChan, hasSubs := broker.CreatePub()
	assert.NotNil(t, subChan)
	assert.NotNil(t, pubChan)
	assert.True(t, hasSubs)
}

func TestMessageReceive(t *testing.T) {
	message := "Good message"
	broker := NewBroker()
	c1 := broker.CreateSub()
	c2 := broker.CreateSub()
	go broker.ReceiveMsg(message)
	assert.Equal(t, message, <-c1)
	assert.Equal(t, message, <-c2)
}

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
