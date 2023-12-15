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

func PubConn(broker pubsub.BrokerConnection, ch chan string) (chan error, error) {
	tlsConfig, quicConfig := createConfigs()
	tr := createTr(pubPort)
	ln, err := tr.Listen(tlsConfig, quicConfig)
	errCh := make(chan error)
	if err != nil {
		return errCh, err
	}

	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				errCh <- err
				break
			}
			go func() {
				for {
					if conn.Context().Err() != nil {
						break
					}
					str, err := conn.AcceptStream(context.Background())
					if err != nil {
						errCh <- err
						continue
					}
					c, hasSubs := broker.CreatePub()
					go func() {
						<-ch
						str.Close()
						conn.CloseWithError(0, "nil")
						broker.Close(c)

					}()
					if !hasSubs {
						str.Write([]byte(pubsub.NoSubsMsg))
					}
					go func() {
						for {
							buf := make([]byte, 512)
							n, err := str.Read(buf)
							if err != nil {
								errCh <- err
								broker.Close(c)
								str.Close()
								break
							}
							broker.ReceiveMsg((string(buf[:n])))
						}
					}()
					go func() {
						for i := range c {
							msg := i
							str.Write([]byte(msg))
						}
					}()
				}
			}()

		}
	}()
	return errCh, nil
}

func SubConn(broker pubsub.BrokerConnection) (chan error, error) {
	tlsConfig, quicConfig := createConfigs()
	tr := createTr(subPort)
	ln, err := tr.Listen(tlsConfig, quicConfig)
	errCh := make(chan error)
	if err != nil {
		return errCh, err
	}

	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				ln.Close()
				errCh <- err
			}
			go func() {
				for {
					str, err := conn.AcceptStream(context.Background())
					if err != nil {
						str.Close()
						errCh <- err
					}
					c := broker.CreateSub()
					go func() {
						for i := range c {
							msg := i
							str.Write([]byte(msg))
						}
					}()
				}
			}()
		}
	}()
	return errCh, nil
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
