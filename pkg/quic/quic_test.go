package quic

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/RarePepeCode/quick-broker/pkg/pubsub"
	"github.com/quic-go/quic-go"
	"github.com/stretchr/testify/assert"
)

var (
	message   = "Good massage"
	pubAdress = serverIp.String() + ":" + strconv.Itoa(pubPort)
	subAdress = serverIp.String() + ":" + strconv.Itoa(subPort)
)

// Tests Pub Connection by creating Client and connecting to it.
// As no Subs are connected to the server, info message about it should be sent via quic.
// By receiving that message we can confirm that we create conenction between server and client.
func TestPubConnectionWithNoSubs(t *testing.T) {
	broker := pubsub.NewBroker()
	PubConn(broker)

	// Mock client connection to the server
	tlsConfig, quicConfig := createConfigs()

	conn, err := quic.DialAddr(context.Background(), pubAdress, tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 512)
	n, err := stream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, pubsub.NoSubsMsg, (string(buf[:n])))

}

// Tests Sub Connection by creating Client and connecting to it.
// Broker Receives the message as it where from the pub and sends to all sub client
func TestSubCreation(t *testing.T) {
	broker := pubsub.NewBroker()
	SubConn(broker)

	tlsConfig, quicConfig := createConfigs()

	// Connect to Sub
	subConn, err := quic.DialAddr(context.Background(), subAdress, tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}

	subStream, err := subConn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = subStream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	// Wait for connection to be establish
	time.Sleep(3 * time.Second)

	// Immitate message publishing to the server
	broker.ReceiveMsg(message)

	buf := make([]byte, 512)
	n, err := subStream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, message, (string(buf[:n])))
}

// First pub connection is establish at that point no subs exist, so server sends message about it.
// After sub connection is made, which will sent the message about it to the first (pub) client.
func TestInformSubConnnected(t *testing.T) {
	broker := pubsub.NewBroker()

	PubConn(broker)
	SubConn(broker)

	tlsConfig, quicConfig := createConfigs()

	// Connect to Pub
	pubConn, err := quic.DialAddr(context.Background(), pubAdress, tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}
	pubStream, err := pubConn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = pubStream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 512)
	n, err := pubStream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// At this point we confirm that no subs are conneted to the server
	assert.Equal(t, pubsub.NoSubsMsg, (string(buf[:n])))

	// Connect to Sub
	subConn, err := quic.DialAddr(context.Background(), subAdress, tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}
	subStream, err := subConn.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = subStream.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}

	buf = make([]byte, 512)
	n, err = pubStream.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	// If sub connected sucssefully pubs client should receive message about it
	assert.Equal(t, pubsub.SubConnectedMsg, (string(buf[:n])))

}
