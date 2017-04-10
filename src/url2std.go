package main

import (
	"os"
	"fmt"
	"encoding/base64"
)
func main() {
	entryURI := os.Args[1]
	str, err := base64.URLEncoding.DecodeString(entryURI)
	if err != nil {
		fmt.Println("err", err.Error())
	}

	fmt.Println("rpc.decode", base64.StdEncoding.EncodeToString(str))
}