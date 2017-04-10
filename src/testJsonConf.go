package main

import (
	"os"
	"fmt"

	"qbox.us/cc/config"
	//"qbox.us/urlrewrite.v2"
)

func main () {
	var table s2
	err := config.LoadFile(&table, os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "config.Load failed:", err)
		os.Exit(-1)
	}

	fmt.Printf("p :%#v\n", table)
}

type slice [] struct {
	A int  `json:"a"`
	B string `json:"b"`
}

type s struct {
	A int
	B string
}

type s2 struct {
	Slice slice `json:"slice"`
        Size int `json:"size"`
}