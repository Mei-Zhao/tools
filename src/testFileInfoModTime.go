package main

import (
	"os"
	"fmt"
)

func main () {
	fileInfo, err := os.Stat("/Users/zhaomei/qbox/tools/src/file/a.txt")
	if err != nil {
		fmt.Println("err", err)
		os.Exit(-5)
	}
	fmt.Println("modeTime", fileInfo.ModTime())
	fmt.Println("modTime", fileInfo.ModTime().UnixNano()/100)
}
