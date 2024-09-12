package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var help = flag.Bool("help", false, "this help")
var isServer = flag.Bool("isServer", false, "is server")
var useTls = flag.Bool("useTls", false, "use tls protocol")
var host = flag.String("host", "", "ip address")
var port = flag.Int("port", 8080, "port")
var password = flag.String("password", "changeme", "password")

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  clipboard_share -host 192.168.1.1 -port 8080
  clipboard_share -host 192.168.1.1 -port 8080 -useTls
  clipboard_share -host 192.168.1.1 -port 8080 -useTls -password 123456
  clipboard_share -host 192.168.1.1 -port 8080 -isServer
  clipboard_share -host 192.168.1.1 -port 8080 -isServer -useTls
  clipboard_share -host 192.168.1.1 -port 8080 -isServer -useTls -password 123456

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	host, err := validateIp(*host)
	if err != nil {
		fmt.Println("Incorrect host")
		return
	}
	port, err := validatePort(*port)
	if err != nil {
		fmt.Println("Incorrect port")
		return
	}
	if *isServer {
		err := runServer(host, port, *useTls, *password)
		if err != nil {
			return
		}
	}
	time.Sleep(1 * time.Second)
	runClient(host, port, *useTls, *password)
}
