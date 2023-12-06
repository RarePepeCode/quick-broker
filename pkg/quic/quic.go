package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

const (
	pubPort = 1234
	subPort = 5678
)

var tlsConfig = &tls.Config{
	Certificates:       []tls.Certificate{loadCert()},
	InsecureSkipVerify: true,
}
var quicConfig = &quic.Config{
	RequireAddressValidation: func(a net.Addr) bool { return false },
}

func PubConn() {
	tr := pubTr()
	ln, err := tr.Listen(tlsConfig, quicConfig)
	if err != nil {
		fmt.Println("Error")
	}
	go func() {
		for {
			conn, err := ln.Accept(context.Background())
			if err != nil {
				fmt.Println("Error")
			}

			go func() {
				for {
					str, err := conn.AcceptStream(context.Background())
					if err != nil {
						fmt.Println(err)
					}
					buf := make([]byte, 512)
					n, err := str.Read(buf)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(string(buf[:n]))
				}
			}()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
	defer cancel()

	conn, err := tr.Dial(ctx, tr.Conn.LocalAddr(), tlsConfig, quicConfig)
	if err != nil {
		fmt.Println(err)
	}
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		fmt.Println(err)
	}
	stream.Write([]byte("Sending Message"))
	time.Sleep(5 * time.Second)
}

func MockClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
	defer cancel()
	tr := pubTr()
	conn, err := tr.Dial(ctx, tr.Conn.LocalAddr(), tlsConfig, quicConfig)
	if err != nil {
		fmt.Println(err)
	}
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		fmt.Println(err)
	}
	stream.Write([]byte("Sending Message"))
	time.Sleep(5 * time.Second)
}

func loadCert() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("cert/cert.pem", "cert/key.pem")
	if err != nil {
		fmt.Println(err)
	}
	return cert
}

func pubTr() quic.Transport {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: pubPort,
	})
	if err != nil {
		fmt.Println(err)
	}
	return quic.Transport{
		Conn: udpConn,
	}
}
