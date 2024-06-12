package main

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"golang.design/x/clipboard"
)

func runClient(ip string, port int) {

	conn, cerr := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
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
			// if string(lastContent) == string(content) {
			// 	continue
			// }
			t.Send("string", content)
			lastContent = content
		case msg := <-watchCh:
			// fmt.Println("get clipboard data: ", string(msg.Type))
			if msg.Type == "string" {
				clipboard.Write(clipboard.FmtText, msg.Content)
			} else if msg.Type == "png" {
				fmt.Println("test png")
				clipboard.Write(clipboard.FmtImage, msg.Content)
			} else {
				fmt.Println("warning: type error")
			}

			lastContent = msg.Content
		case content := <-imageCh:
			if bytes.Equal(lastContent, content) {
				continue
			}
			t.Send("png", content)
			lastContent = content
		}

	}

}
