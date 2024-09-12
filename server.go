package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
)

var lock sync.Mutex
var tcpList map[string]*Tcp

const (
	certFilePath = "cert/cert.pem"
	keyFilePath  = "cert/key.pem"
)

func runServer(host string, port int, useTls bool, password string) error {
	var ln net.Listener
	var err error
	if useTls {
		cert, certErr := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		if certErr != nil {
			panic(certErr)
		}
		config := tls.Config{Certificates: []tls.Certificate{cert}}
		ln, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", host, port), &config)
	} else {
		ln, err = net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	}
	if err != nil {
		panic(err)
	}
	tcpList = make(map[string]*Tcp)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				fmt.Println("Accep error", err.Error())
				break
			}
			t := NewTcp(c)
			if !passwordIsRighted(t, password) {
				t.Send(TMSystem, (&TMSystemMsg{403, []byte("Password is incorrect!")}).Bytes())
				fmt.Println("host: ", t.conn.RemoteAddr().String(), " password error.")
				t.Close()
				continue
			} else {
				t.Send(TMSystem, (&TMSystemMsg{200, []byte("Password is righted!")}).Bytes())
			}
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
		fmt.Printf("Get the msg type  %d and notify anyone.\n", msg.Type)
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

func passwordIsRighted(t *Tcp, password string) bool {
	msg, error := t.Read()
	if error != nil {
		return false
	}
	if msg.Type != TMPassword || string(msg.Content) != password {
		return false
	}
	return true
}
