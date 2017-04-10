package  main


import (
	"encoding/base64"
	"os"
	"fmt"
)

func main() {
	str := os.Args[1]
	bytes, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(string(bytes))
}
