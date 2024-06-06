package main

import (
	"context"
	"fmt"

	"golang.design/x/clipboard"
)

func main() {
	fmt.Println("init cliboard...")
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	var content []byte
	for {
		content = <-ch
		fmt.Println("w: ", string(content))
	}

	// bstr := clipboard.Read(clipboard.FmtText)
	// fmt.Println(string(bstr))
}
