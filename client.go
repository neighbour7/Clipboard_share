package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"

	"golang.design/x/clipboard"
)

func runClient(host string, port int, useTls bool, password string) {
	var conn net.Conn
	var err error
	if useTls {
		config := tls.Config{InsecureSkipVerify: true}
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &config)
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	}
	if err != nil {
		panic(err)
	}

	t := NewTcp(conn)
	if !validatePassword(t, password) {
		fmt.Println("password error!")
		t.Close()
		return
	}
	watchCh := t.Watch()
	defer t.Close()
	err = clipboard.Init()
	if err != nil {
		panic(err)
	}
	textCh := clipboard.Watch(context.TODO(), clipboard.FmtText)
	imageCh := clipboard.Watch(context.Background(), clipboard.FmtImage)
	var lastContent []byte
	var content []byte

	for {
		select {
		case content = <-textCh:
			if bytes.Equal(lastContent, content) {
				continue
			}
			t.Send(TMText, content)
			fmt.Println("Send text: ", string(content))
			lastContent = content
		case msg := <-watchCh:
			if msg.Type == TMText {
				clipboard.Write(clipboard.FmtText, msg.Content)
				fmt.Println("Write text: ", string(msg.Content))
			} else if msg.Type == TMImg {
				clipboard.Write(clipboard.FmtImage, msg.Content)
				fmt.Println("Write png size: ", len(msg.Content))
			} else {
				fmt.Println("Warning: type error")
			}

			lastContent = msg.Content
		case content := <-imageCh:
			if bytes.Equal(lastContent, content) {
				continue
			}
			t.Send(TMImg, content)
			fmt.Println("Send png size: ", len(content))
			lastContent = content
		}
	}
}

func validatePassword(t *Tcp, password string) bool {
	error := t.Send(TMPassword, []byte(password))
	if error != nil {
		return false
	}
	msg, error := t.Read()

	if msg.Type == TMSystem {
		ts := &TMSystemMsg{}
		err := json.Unmarshal(msg.Content, ts)
		if err != nil {
			return false
		}
		if ts.Type == 200 {
			return true
		}
		return false
	}
	return false
}
