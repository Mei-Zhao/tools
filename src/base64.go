package main

import (
	"encoding/base64"
	"os"
	"fmt"
)

func main() {
	str := os.Args[1] + ":" + os.Args[2]
	fmt.Println("encodeEntry", base64.URLEncoding.EncodeToString([]byte(str)))
}
