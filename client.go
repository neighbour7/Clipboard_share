package main

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"golang.design/x/clipboard"
)

func runClient(host string, port int) {

	conn, cerr := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if cerr != nil {
		fmt.Println("Error connecting:", cerr.Error())
		return
	}
	t := NewTcp(conn)
	watchCh := t.Watch()
	defer t.Close()
	err := clipboard.Init()
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
			t.Send("string", content)
			fmt.Println("Send text: ", string(content))
			lastContent = content
		case msg := <-watchCh:
			if msg.Type == "string" {
				clipboard.Write(clipboard.FmtText, msg.Content)
				fmt.Println("Write text: ", string(msg.Content))
			} else if msg.Type == "png" {
				clipboard.Write(clipboard.FmtImage, msg.Content)
				fmt.Println("Write png size: ", len(msg.Content))
			} else {
				fmt.Println("warning: type error")
			}

			lastContent = msg.Content
		case content := <-imageCh:
			if bytes.Equal(lastContent, content) {
				continue
			}
			t.Send("png", content)
			fmt.Println("Send png size: ", len(content))
			lastContent = content
		}
	}
}
