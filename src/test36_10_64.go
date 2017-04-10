package main

import (
	"strconv"
	"fmt"
)

func main () {
	i := 50
	fileNname := strconv.FormatInt(int64(i),36)
	fmt.Println(fileNname)
}
