package main

import (
	"os"
	"fmt"
	"strconv"
)

func main() {
	psect := os.Args[1]
	n, err := strconv.ParseUint(psect,10,64)
	if err != nil {
		fmt.Println("strconv.ParseUint", psect)
	}
	diskId, isSector := DecodePsect(n)
	fmt.Println("diskId", diskId)
	fmt.Println("isSector", isSector)
}

func DecodePsect(psect uint64) (diskId, iSector uint32) {
	diskId, iSector = uint32(psect>>32), uint32(psect)
	return
}
