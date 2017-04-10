package main

import (
	"fmt"
	"time"

)


func main(){

	fmt.Println(time.Unix(time.Now().UnixNano()/1e9, 0).Format(time.RFC1123Z))
	fmt.Println(time.Unix(time.Now().UnixNano()/1e9, 0).Format("02/Jan/2006:15:04:05 -0700"))
	timestamp := "2017-02-10T00:59:49.457753Z"
	time, err := time.Parse(time.RFC3339,timestamp)
	if err != nil {
		fmt.Println(time)
	}

	fmt.Println(time.Format("02/Jan/2006:15:04:05 -0700"))
}