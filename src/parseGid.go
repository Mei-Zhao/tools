package main

import (
	"encoding/binary"
	"log"
	"os"
	"qbox.us/pfd/api/types"
	"fmt"
	"time"
)

type Gid [12]byte

func (self Gid) MotherDgid() uint32 {
	return binary.LittleEndian.Uint32(self[:4])
}

func main() {
	egid := os.Args[1]
	gid, err := types.DecodeGid(egid)
	if err != nil {
		log.Println(egid)
	}
	fmt.Println("dgid",gid.MotherDgid())
	fmt.Println(egid, time.Unix(gid.UnixTime()/1e9, gid.UnixTime()%1e9).Format("2006-01-02_15:04"))
}
