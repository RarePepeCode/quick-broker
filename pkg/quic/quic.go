package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/RarePepeCode/quick-broker/pkg/pubsub"
	"github.com/quic-go/quic-go"
)

var (
	pubPort = 1234
	subPort = 5678
)

func PubConn(broker pubsub.BrokerConnection) {
	tlsConfig, quicConfig := createConfigs()
	tr := pubTr(pubPort)
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
					stremComm(broker, str, c)
				}
			}()
		}
	}()
}

func SubConn(broker pubsub.BrokerConnection) {
	tlsConfig, quicConfig := createConfigs()
	tr := pubTr(subPort)
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
					stremComm(broker, str, c)
				}
			}()
		}
	}()
}

func ClientMain() error {
	message := "Goog massage"
	tlsConfig, quicConfig := createConfigs()

	conn, err := quic.DialAddr(context.Background(), "127.0.0.1:1234", tlsConfig, quicConfig)
	if err != nil {
		return err
	}
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	_, err = stream.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func stremComm(broker pubsub.BrokerConnection, str quic.Stream, c chan string) {
	go func() {
		for {
			buf := make([]byte, 512)
			n, err := str.Read(buf)
			if err != nil {
				fmt.Println(err)
			}
			broker.ReceiveMsg((string(buf[:n])))
		}
	}()
	go func() {
		for {
			msg := <-c
			fmt.Println("Sending to client - ", msg)
			str.Write([]byte(msg))
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
	}
	return tlsConfig, quicConfig
}

func pubTr(port int) quic.Transport {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	})
	if err != nil {
		fmt.Println(err)
	}
	return quic.Transport{
		Conn: udpConn,
	}
}
