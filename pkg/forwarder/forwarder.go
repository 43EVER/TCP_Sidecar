package forwarder

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type ProxyPairItem struct {
	Conn *net.TCPConn
	Name string
}

func (reader *ProxyPairItem) CopyTo(writer *ProxyPairItem) error {
	_, err := io.Copy(reader.Conn, writer.Conn)
	switch err {
	case nil:
		// go closeHelper(writer.Conn.CloseWrite, fmt.Sprintf("Closing Write[%s]", writer.Name))
		closeHelper(reader.Conn.CloseRead, fmt.Sprintf("Closing Read[%s]", reader.Name))
		return nil
	default:
		return err
	}
}

type ProxyPair struct {
	PassiveItem ProxyPairItem
	ActiveItem  ProxyPairItem
	Name        string
}

func closeHelper(fun func() error, errorPrefix string) {
	err := fun()
	if err != nil {
		log.Printf("%s, error: %s", errorPrefix, err.Error())
	}
}

func copyHelper(readItem *ProxyPairItem, writeItem *ProxyPairItem, wg *sync.WaitGroup) {
	defer wg.Done()

	err := readItem.CopyTo(writeItem)
	if err != nil {
		fmt.Printf("[%s] Copy to [%s], error: %v", readItem.Name, writeItem.Name, err)
	}
}

func (p *ProxyPair) StartForward() {
	log.Printf("'%s' start forwarding '%s' <-> '%s'",
		p.Name, p.ActiveItem.Name, p.PassiveItem.Name)
	defer closeHelper(p.PassiveItem.Conn.CloseWrite, fmt.Sprintf("Closing write [%s]", p.PassiveItem.Name))
	defer closeHelper(p.ActiveItem.Conn.CloseWrite, fmt.Sprintf("Closing  write [%s]", p.ActiveItem.Name))

	// defer closeHelper(p.PassiveItem.Conn.Close, fmt.Sprintf("Closing [%s]", p.PassiveItem.Name))
	// defer closeHelper(p.ActiveItem.Conn.Close, fmt.Sprintf("Closing  [%s]", p.ActiveItem.Name))

	var wg sync.WaitGroup
	wg.Add(2)
	go copyHelper(&p.ActiveItem, &p.PassiveItem, &wg)
	go copyHelper(&p.PassiveItem, &p.ActiveItem, &wg)
	wg.Wait()

	log.Printf("'%s' finish forwarding '%s' <-> '%s'",
		p.Name, p.ActiveItem.Name, p.PassiveItem.Name)
}
