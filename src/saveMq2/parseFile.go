package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	//"time"

	//"qbox.us/qmq/qmqapi/v1/mq"

	//"github.com/ngaut/log"
	//"github.com/qiniu/api/auth/digest"
	xlog "github.com/qiniu/xlog.v1"

	"errors"
	"strconv"
	"strings"

	"qbox.us/cc/table"
	//"fmt"
	"fmt"
	"log"
	pfdtypes "qbox.us/pfd/api/types"
)

func main() {
	//mqHost := flag.String("mqhost", "http://192.168.58.61:14500", "mq host")
	//flag.Parse()
	//trans := digest.NewTransport(&digest.Mac{
	//	AccessKey: "-s_NZdOCsmL9vRsY_TAhxnBuDLh4dgh-HqM8eRV_",
	//	SecretKey: []byte("3b2ghfzvvfP8laMGb9z5c_zulNTZ1pzbHa008wDZ"),
	//}, nil)
	//mqcli := mq.NewWithHost(trans, *mqHost)
	decoder := json.NewDecoder(os.Stdin)
	save := flag.String("save", "save_parse", "save file")
	flag.Parse()
	var msg []byte
	xl := xlog.NewDummy()

	saveFile, err := os.Create(*save)
	if err != nil {
		log.Fatal(err)
	}
	defer saveFile.Close()
	encoder := json.NewEncoder(saveFile)
	for {
		err := decoder.Decode(&msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			xl.Fatal(err)
		}
		//success := false
		//for i := 0; i < 20; i++ {
		//	_, err = mqcli.Put(xl, "bddelete", msg)
		//	if err == nil {
		//		success = true
		//		//fmt.Println("true")
		//		break
		//	}
		//	xl.Warn(err)
		//	time.Sleep(time.Second)
		//}
		//if !success {
		//	xl.Fatal(err)
		//}
		_, err = parseDeleteMessage(msg)
		if err != nil {
			log.Fatal(err)
		}
		err = encoder.Encode(msg)
		if err != nil {
			log.Fatal(err)
		}

	}
}

const (
	FhPfdLen = pfdtypes.FileHandleBytes
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

func parseDeleteMessage(msg []byte) (m *deleteMessage, err error) {
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
	fields := strings.Split(string(msg[FhPfdLen:]), "\t")

	if len(fields) < 6 {
		fmt.Println("119 msg:", string(msg))
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
