package main

import (
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
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	var content []byte
	for {
		select {
		case content = <-ch:
			t.Send("string", content)
		case msg := <-watchCh:
			fmt.Println(string(msg.Content))
			clipboard.Write(clipboard.FmtText, msg.Content)
		}

	}

}
