package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	i, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Unix(i/1e7, i%1e7*100))
}