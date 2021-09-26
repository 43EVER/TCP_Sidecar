package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	forwarder "tcpsidecar/pkg/forwarder"
)

type ID struct {
	lock sync.Mutex
	id   int
}

func (id *ID) Get() int {
	id.lock.Lock()
	defer id.lock.Unlock()

	id.id += 1

	return id.id
}

func (id *ID) Cancel() {
	id.lock.Lock()
	defer id.lock.Unlock()

	id.id -= 1
}

func proxyStart(listenAddr, targetAddr string) {
	listenTCPAddr, err := net.ResolveTCPAddr("tcp4", listenAddr)
	if err != nil {
		panic("listenAddr resolve failed")
	}

	listener, err := net.ListenTCP("tcp4", listenTCPAddr)
	if err != nil {
		if err != nil {
			fmt.Printf("Unable to listen on: %s, error: %s\n", listenAddr, err.Error())
			os.Exit(1)
		}
	}

	id := ID{}

	for {
		passiveConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Unable to accept a connection, error: %s\n", err.Error())
			continue
		}

		go func() {
			targetTCPAddr, err := net.ResolveTCPAddr("tcp4", targetAddr)
			if err != nil {
				panic("targetAddr resolve failed")
			}

			activeConn, err := net.DialTCP("tcp4", nil, targetTCPAddr)
			if err != nil {
				fmt.Printf("Unable to connect to: %s, error: %s\n", targetAddr, err.Error())
			}

			proxyPair := forwarder.ProxyPair{
				ActiveItem: forwarder.ProxyPairItem{
					Name: fmt.Sprintf("activeConn(%s<->%s)", activeConn.RemoteAddr().String(), activeConn.LocalAddr().String()),
					Conn: activeConn,
				},
				PassiveItem: forwarder.ProxyPairItem{
					Name: fmt.Sprintf("passiveConn(%s<->%s)", passiveConn.RemoteAddr().String(), passiveConn.LocalAddr().String()),
					Conn: passiveConn,
				},
				Name: fmt.Sprintf("proxyPair(%d)", id.Get()),
			}

			proxyPair.StartForward()
			id.Cancel()
		}()
	}
}

func main() {
	var listenAddr string
	var targetAddr string

	flag.StringVar(&listenAddr, "l", ":80", "监听地址， 默认为 :80")
	flag.StringVar(&targetAddr, "t", ":8080", "转发地址，默认为 :8080")
	flag.Parse()

	proxyStart(listenAddr, targetAddr)
}
