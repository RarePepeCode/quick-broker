package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/RarePepeCode/quick-broker/pkg/pubsub"
	"github.com/quic-go/quic-go"
)

var (
	serverIp = net.IPv4(127, 0, 0, 1)
	pubPort  = 1234
	subPort  = 5678
)

func PubConn(broker pubsub.BrokerConnection) {
	tlsConfig, quicConfig := createConfigs()
	tr := createTr(pubPort)
	ln, err := tr.Listen(tlsConfig, quicConfig)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				fmt.Println(err)
			}
			go func() {
				for {
					str, err := conn.AcceptStream(context.Background())
					if err != nil {
						fmt.Println(err)
					}
					c, hasSubs := broker.CreatePub()
					if !hasSubs {
						fmt.Println(pubsub.NoSubsMsg)
						str.Write([]byte(pubsub.NoSubsMsg))
					}
					go func() {
						for {
							buf := make([]byte, 512)
							n, err := str.Read(buf)
							if err != nil {
								fmt.Println(err)
								broker.Close(c)
								break
							}
							broker.ReceiveMsg((string(buf[:n])))
						}
					}()
					go func() {
						for i := range c {
							msg := i
							fmt.Println("Sending to Pub client - ", msg)
							str.Write([]byte(msg))
						}
					}()
				}
			}()
		}
	}()
}

func SubConn(broker pubsub.BrokerConnection) {
	tlsConfig, quicConfig := createConfigs()
	tr := createTr(subPort)
	ln, err := tr.Listen(tlsConfig, quicConfig)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				fmt.Println(err)
			}
			go func() {
				for {
					str, err := conn.AcceptStream(context.Background())
					if err != nil {
						fmt.Println(err)
					}
					c := broker.CreateSub()
					go func() {
						for i := range c {
							msg := i
							fmt.Println("Sending to Sub client - ", msg)
							str.Write([]byte(msg))
						}
					}()
				}
			}()
		}
	}()
}

func loadCert() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("cert/cert.pem", "cert/key.pem")
	if err != nil {
		fmt.Println(err)
	}
	return cert
}

func createConfigs() (*tls.Config, *quic.Config) {
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{loadCert()},
		InsecureSkipVerify: true,
	}
	quicConfig := &quic.Config{
		RequireAddressValidation: func(a net.Addr) bool { return false },
		MaxIdleTimeout:           1 * time.Hour,
	}
	return tlsConfig, quicConfig
}

func createTr(port int) quic.Transport {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   serverIp,
		Port: port,
	})
	if err != nil {
		fmt.Println(err)
	}
	return quic.Transport{
		Conn: udpConn,
	}
}
