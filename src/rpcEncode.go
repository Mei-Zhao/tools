package main

import (
	"qbox.us/rpc"
	"os"
	"fmt"
)
func main() {
	entryURI := os.Args[1] + ":" + os.Args[2]
	fmt.Println("rpc.Encode", rpc.EncodeURI(entryURI))
}