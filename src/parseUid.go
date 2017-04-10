package main

import (
	"strconv"
	"fmt"
	"os"
)

func main() {
	owner := os.Args[1]
	uint64_owner, err := strconv.ParseUint(owner,10,64)
	if err != nil  {
		fmt.Println("err", err)
	}
	fmt.Println(strconv.FormatUint(uint64_owner, 36))
}
