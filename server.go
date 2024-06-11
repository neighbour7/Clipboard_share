package main

import (
	"fmt"
	"net"
)

var tcpList []*Tcp

// var lock sync.Mutex

func runServer(ip string, port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("list error", err.Error())
		return
	}
	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("accep error", err.Error())
			return
		}
		t := &Tcp{
			conn:    c,
			name:    c.RemoteAddr().String(),
			watchCh: make(chan *TcpMsg, 1),
		}
		tcpList = append(tcpList, t)
		go listenMsgHandler(t)
	}
}

func listenMsgHandler(t *Tcp) {
	for {
		msg, err := t.Read()
		if err != nil {
			break
		}
		notifyMsg(msg)
	}
}

func notifyMsg(content *TcpMsg) {
	for _, t := range tcpList {
		t.Send(content.Type, content.Content)
	}
}
