package main

import (
	"os"
	"fmt"
)

func main() {
	name := "test.txt"
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
	if err != nil {
		fmt.Println("err", err)
	}
	f.WriteString()
}