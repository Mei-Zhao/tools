package main

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"os"
	"sync"

	log "qiniupkg.com/x/log.v7"

	"errors"
	"qbox.us/cc"
	"qbox.us/largefile"
	"strconv"
	"strings"

	"fmt"
	"qbox.us/cc/table"
	pfdType "qbox.us/pfd/api/types"
)

type recordHeader struct {
	MsgLen   uint16
	Tag      uint8
	State    uint8
	LockTime uint32 // in s
}

const (
	indexSize     = 32
	headerSize    = 8
	recordTag     = 0x0a
	offsetOfState = 3
)

const (
	stateNormal     = 0x01
	stateProcessing = 0x02
	stateDone       = 0x04

	stateTravelPut     = stateNormal | stateProcessing | stateDone
	stateTravelGet     = stateProcessing | stateDone
	stateTravelTimeout = stateDone
)

const (
	FhPfdLen = pfdType.FileHandleBytes
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

type Message struct {
	ModifyTime int64
	Uid        uint32
	MsgId      string
	OwnerId    string
	Bucket     string
	Key        string
	Fh         []byte
	Op         string
	Domains    []string
	PutTime    int64
	Itbl       uint32
}

var (
	EInvalidMsgFormat = errors.New("Invalid Message Format")
)

//fields := []string{
//	strconv.FormatUint(uint64(changeTime), 10),
//	strconv.FormatUint(uint64(owner), 10),
//	table.Escape(bucket),
//	table.Escape(key),
//	op,
//	base64.URLEncoding.EncodeToString(fh),
//	strconv.FormatUint(uint64(putTime), 10),
//	strconv.FormatUint(uint64(itbl), 10),
//}
//msg := strings.Join(fields, "\t")

func ParseMessage(id, message string) (m *Message, err error) {
	fields := strings.Split(message, "\t")
	if len(fields) < 4 {
		return nil, EInvalidMsgFormat
	}
	modifyTime, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil, err
	}
	ownerId := fields[1]
	bucket, err := table.Unescape(fields[2])
	if err != nil {
		return nil, err
	}
	key, err := table.Unescape(fields[3])
	if err != nil {
		return nil, err
	}
	uid, err := strconv.ParseUint(ownerId, 10, 64)
	if err != nil {
		return nil, err
	}

	var rsop string
	var fh []byte
	if len(fields) > 5 {
		rsop = fields[4]
		fh, err = base64.URLEncoding.DecodeString(fields[5])
		if err != nil {
			return nil, err
		}
	}

	var putTime int64
	if len(fields) > 6 {
		putTime, err = strconv.ParseInt(fields[6], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	var itbl uint32
	if len(fields) > 7 {
		itbl64, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return nil, err
		}
		itbl = uint32(itbl64)
	}
	return &Message{modifyTime, uint32(uid), id, ownerId, bucket, key, fh, rsop, nil, putTime, itbl}, nil
}

func Delete(m *Message) string {
	fields := []string{
		string(m.Fh),
		strconv.FormatUint(uint64(m.Uid), 10),
		table.Escape(m.Bucket),
		table.Escape(m.Key),
		strconv.FormatInt(m.ModifyTime, 10),
		strconv.FormatInt(m.PutTime, 10),
		table.Escape(m.Op),
	}
	msg := strings.Join(fields, "\t")
	return msg
}

func filterMsg2(wg *sync.WaitGroup, in, out chan []byte) {
	defer wg.Done()
	for msg := range in {

		//读取 parse rsedit的消息
		m, err := ParseMessage("", string(msg))
		if err != nil {
			log.Fatal(err)
		}
		//bddelete过滤
		//op过滤
		if !hasString(deleteFileOps, m.Op) {
			continue
		}
		//bddelete消息内容
		deleteMsg := Delete(m)
		out <- []byte(deleteMsg)
	}
}

var deleteFileOps = []string{
	"delete",
	"ins",
	"put",
	"failed",
}

func hasString(ss []string, s string) bool {

	for _, s0 := range ss {
		if s0 == s {
			return true
		}
	}
	return false
}

func main() {
	startOffset := flag.Int64("s", 0, "start offset")
	limit := flag.Int64("c", -1, "read max count")
	dir := flag.String("d", "", "dir")
	chunkBit := flag.Uint("b", 29, "chunk bit")
	save := flag.String("save", "save", "save file")
	process := flag.Int("p", 5, "process to parse")
	flag.Parse()
	f, err := largefile.Open(*dir, *chunkBit)
	if err != nil {
		log.Fatal(err)
	}
	r := &cc.Reader{f, *startOffset}

	// find real start offset
	var tmp = make([]byte, 1024*512)
	fmt.Println("tmp", 1024*512)
	_, err = io.ReadFull(r, tmp)
	if err != nil {
		log.Fatal(err)
	}
	var realStartOffset, i int64

	//猜测rsedit消息开始的正确位置
	for ; i < 524228; i++ {
		fmt.Println("i: ", i)
		//parse header msgLen , parseMsg
		msgLen := tmp[i : i+2]
		len, err := strconv.ParseInt(string(msgLen),10, 64)
		if err != nil {
			fmt.Println("err: ", err)
		}
		fmt.Println("len of msg: ", len)
		_, err = ParseMessage("", string(tmp[i+2:i+2+len]))
		fmt.Println("err:", err)
		if err == nil {
			realStartOffset = *startOffset + i - 8
			log.Println("real start offset", realStartOffset)
			break
		}

	}

	r = &cc.Reader{f, realStartOffset}
	br := bufio.NewReaderSize(r, 4*1024*1024)
	saveFile, err := os.Create(*save)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("create", *save)
	defer saveFile.Close()
	encoder := json.NewEncoder(saveFile)
	var readCh = make(chan []byte, 1)
	var saveCh = make(chan []byte)
	var wg = &sync.WaitGroup{}
	wg.Add(*process)
	for i := 0; i < *process; i++ {
		go filterMsg2(wg, readCh, saveCh)
	}
	go func() {
		for msg := range saveCh {
			err = encoder.Encode(msg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	for count := int64(0); count != *limit; count++ {
		var h recordHeader
		err = binary.Read(br, binary.LittleEndian, &h)
		if err != nil {
			if err == io.EOF {
				close(readCh)
				break
			}
			log.Fatal(err)
		}
		msg := make([]byte, h.MsgLen)
		n, err := io.ReadFull(br, msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(n, err)
		}
		readCh <- msg
	}
	wg.Wait()
	close(saveCh)
}
