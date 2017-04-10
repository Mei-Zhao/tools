package main

import (
	"os"
	"fmt"
	"encoding/base64"
)
func main() {
	str := os.Args[1]
	fmt.Println("encodeEntry", base64.URLEncoding.EncodeToString([]byte(str)))
}

