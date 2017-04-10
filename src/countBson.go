package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

func readInt32(rd io.Reader) (int32, error) {
	var i int32
	if err := binary.Read(rd, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var count uint64
	for {
		n, err := readInt32(f)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		_, err = f.Seek(int64(n-4), 1)
		if err != nil {
			log.Fatal(err)
		}
		count += 1
		if count%100000000 == 0 {
			log.Println(count)
		}
	}
	log.Println("count:", count)
}
