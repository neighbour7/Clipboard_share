package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args
	ip, err := validateIp(args[1])
	if err != nil {
		fmt.Println("IP地址错误")
		return
	}
	port, err := validatePort(args[2])
	if err != nil {
		fmt.Println("端口错误")
		return
	}
	if args[3] == "server" {
		err := runServer(ip, port)
		if err != nil {
			return
		}
	}
	time.Sleep(1 * time.Second)
	runClient(ip, port)
}
