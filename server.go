package main

import (
	"fmt"
	"net"
	"sync"
)

var lock sync.Mutex
var tcpList map[string]*Tcp

func runServer(ip string, port int) error {
	tcpList = make(map[string]*Tcp)
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("list error", err.Error())
		return err
	}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				fmt.Println("accep error", err.Error())
				break
			}
			t := NewTcp(c)
			lock.Lock()
			tcpList[c.RemoteAddr().String()] = t
			lock.Unlock()
			fmt.Println("new conn: ", c.RemoteAddr().String())
			go listenMsgHandler(t)
		}
	}()
	return nil
}

func listenMsgHandler(t *Tcp) {
	defer t.Close()
	for {
		msg, err := t.Read()
		if err != nil {
			break
		}
		fmt.Printf("Get the %s and notify anyone.", msg.Type)
		notifyMsg(msg)
	}
	lock.Lock()
	delete(tcpList, t.name)
	lock.Unlock()
}

func notifyMsg(content *TcpMsg) {
	for _, t := range tcpList {
		t.Send(content.Type, content.Content)
	}
}
