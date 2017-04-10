package main

import (
	"strconv"
	"fmt"
)

func main () {
	owner := uint32(1377015735)
	fmt.Println(strconv.FormatUint(uint64(owner), 36))
}
