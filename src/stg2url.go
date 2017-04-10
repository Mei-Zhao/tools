package main

import (
	"os"
	"fmt"
	"encoding/base64"
)
func main() {
	entryURI := os.Args[1]
	str, err := base64.StdEncoding.DecodeString(entryURI)
	if err != nil {
		fmt.Println("err", err.Error())
	}

	fmt.Println("rpc.decode", base64.URLEncoding.EncodeToString(str))
}