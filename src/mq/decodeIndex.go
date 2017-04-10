package main

import (
	"encoding/binary"
	"flag"
	"log"
	"os"
)

type indexFile struct {
	OffPut     int64
	OffGet     int64
	OffTimeout int64
	Expires    uint32
	Reserved   uint32
}

func main() {
	expires := flag.Uint64("o", 0, "offset timeout")
	flag.Parse()
	if *expires == 0 {
		log.Fatalln("offTimeout not set", *expires)
	}
	expires2:= uint32(*expires)
	log.Println(expires2)
	findex, err := os.Open("INDEX")
	if err != nil {
		log.Fatal(err)
	}
	var index indexFile
	err = binary.Read(findex, binary.LittleEndian, &index)
	if err != nil {
		log.Println(err)
	}
	findex.Close()
	log.Printf("old index %+v", index)
	index.Expires = expires2
	findexW, err := os.Create("INDEX")
	if err != nil {
		log.Fatalln(err)
	}
	err = binary.Write(findexW, binary.LittleEndian, index)
	if err != nil {
		log.Fatalln(err)
	}
	findexW.Close()
}
