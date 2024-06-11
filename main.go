package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	ip, err := validateIp(args[2])
	if err != nil {
		fmt.Println("IP地址错误")
		return
	}
	port, err := validatePort(args[3])
	if err != nil {
		fmt.Println("端口错误")
		return
	}
	if args[1] == "client" {
		runClient(ip, port)
	} else if args[1] == "server" {
		runServer(ip, port)
	} else {
		fmt.Println("参数错误！")
	}
}
