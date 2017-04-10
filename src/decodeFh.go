package main

import (
	"qbox.us/fh/fhver"
	"github.com/qiniu/encoding/binary"

	"fmt"
	"bytes"
	"syscall"
	"encoding/base64"
	"os"
)

const (
	FileHandle_ChunkBits = 0x96

	FileHandleBytes = 1 + 1 + 2 + 12 + 8 + 8 + 8 + 20
)

type Gid [12]byte

func EncodeGid(gid Gid) string {
	return base64.URLEncoding.EncodeToString(gid[:])
}

type FileHandle struct {
	Ver    uint8    // 1  Byte
	Tag    uint8    // 1  Byte
	Upibd  uint16   // 2  Byte
	Gid    Gid      // 12 Byte
	Offset int64    // 8  Byte
	Fsize  int64    // 8  Byte
	Fid    uint64   // 8  Byte
	Hash   [20]byte // 20 Byte
}

func main() {
	fh := os.Args[1]
	fhByte, _ := base64.StdEncoding.DecodeString(fh)
	fhi, _:= DecodeFh(fhByte)
	fmt.Printf("decodeFh : %#v\n",fhi)
	fmt.Println("gid",EncodeGid(fhi.Gid) )
	fmt.Println("size", fhi.Fsize)
	fmt.Println("offset", fhi.Offset)
	fmt.Println("fid", fhi.Fid)
}

func DecodeFh(fh []byte) (fhi *FileHandle, err error) {

	if len(fh) != FileHandleBytes ||
	(fh[0] != fhver.FhPfd && fh[0] != fhver.FhPfdV2) ||
	fh[1] != FileHandle_ChunkBits {

		err = syscall.EINVAL
		return
	}
	fhi = new(FileHandle)

	r := bytes.NewReader(fh)
	err = binary.Read(r, binary.LittleEndian, fhi)
	return
}