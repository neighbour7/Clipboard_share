package main

import (
	"fmt"
	"net"
	"sync"
)

var lock sync.Mutex
var tcpList map[string]*Tcp

func runServer(host string, port int) error {
	tcpList = make(map[string]*Tcp)
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				fmt.Println("Accep error", err.Error())
				break
			}
			t := NewTcp(c)
			lock.Lock()
			tcpList[c.RemoteAddr().String()] = t
			lock.Unlock()
			fmt.Println("Add host: ", c.RemoteAddr().String())
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
		fmt.Printf("Get the %s and notify anyone.\n", msg.Type)
		notifyMsg(msg)
	}
	lock.Lock()
	delete(tcpList, t.name)
	fmt.Printf("Delete host %s\n", t.name)
	lock.Unlock()
}

func notifyMsg(content *TcpMsg) {
	for _, t := range tcpList {
		t.Send(content.Type, content.Content)
	}
}
