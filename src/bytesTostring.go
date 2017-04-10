package main

import (
	"fmt"
	"errors"
	"strings"
	"strconv"

	"qbox.us/cc/table"
	"github.com/qiniu/log.v1"
	pfdtypes "qbox.us/pfd/api/types"
)

type deleteMessage struct {
	ver        int
	fh         []byte
	uid        uint32
	bucket     string
	key        string
	modifyTime int64
	putTime    int64
	op         string
}

func main() {
	//var content []byte
	//content := {7 150 0 236 96 33 202 154 159 89 30 228 6 20 97 149 1 0 0 0 0 0 120 145 181 41 18 0 0 0 6 0 22 135 186 207 59 201 159 195 16 107 41 21 164 56 97 43 244 83 22 74 188 9 49 51 50 51 51 53 57 49 51 52 9 99 100 110 46 51 54 48 105 110 46 99 111 109 9 97 102 57 48 100 102 49 52 51 98 56 56 53 56 100 97 102 57 49 48 98 51 49 99 98 52 55 49 53 49 100 52 9 49 52 57 48 52 54 53 50 57 51 9 49 51 53 50 56 54 57 51 48 50 9 100 101 108 101 116 101}
	content := []byte{7,150,0,236,96,33,202,154,159,89,30,228,6,20,35,171,4,0,0,0,0,0,138,49,186,51,18,0,0,0,6,0,22,199,57,215,201,56,174,74,177,58,177,176,46,253,16,46,7,49,28,55,203,9,49,51,50,51,51,53,57,49,51,52,9,99,100,110,46,51,54,48,105,110,46,99,111,109,9,101,48,51,100,50,99,97,50,49,101,51,99,52,53,97,101,50,54,49,100,54,55,97,56,102,53,101,97,101,55,51,53,9,49,52,57,48,52,54,53,50,55,57,9,49,51,53,51,49,53,53,54,56,50,9,100,101,108,101,116,101}
	fmt.Println(string(content))
	m ,err := parseDeleteMessage2(content)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("%+v\n", m)

}

const (
	FhPfdLen = pfdtypes.FileHandleBytes
)


func parseDeleteMessage2(msg []byte) (m *deleteMessage, err error) {
	if len(msg) < FhPfdLen {
		err = errors.New("msg is too short, invalid fh")
		return
	}

	if len(msg) == FhPfdLen {
		// for compatibility when migration, fh can only be FhPfd or FhPfdV2
		m = &deleteMessage{
			ver: 0,
			fh:  msg,
		}
		return
	}

	var info deleteMessage
	info.ver = 1

	// has modifyTime & putTime:   fh \t uid \t bucket \t key \t modifyTime \t putTime \t op
	info.fh = msg[:FhPfdLen]
	fmt.Printf("fh: %+v\n", info.fh)
	fields := strings.Split(string(msg[FhPfdLen:]), "\t")
	fmt.Printf("fields:%+v\n", fields)

	if len(fields) < 6 {
		err = errors.New("invalid message")
		return
	}

	// fields[0] is empty
	var uid uint64
	uid, err = strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return
	}
	info.uid = uint32(uid)

	info.bucket, err = table.Unescape(fields[2])
	if err != nil {
		return
	}
	info.key, err = table.Unescape(fields[3])
	if err != nil {
		return
	}

	info.modifyTime, err = strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		return
	}
	info.putTime, err = strconv.ParseInt(fields[5], 10, 64)
	if err != nil {
		return
	}
	if len(fields) >= 7 {
		info.op, err = table.Unescape(fields[6])
		if err != nil {
			return
		}
	}

	m = &info
	return
}
