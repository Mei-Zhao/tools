package main

import (
	"qbox.us/rpc"
	"os"
	"fmt"
)
func main() {
	encodeURI := os.Args[1]
	u, err := rpc.DecodeURI(encodeURI)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(u)
}