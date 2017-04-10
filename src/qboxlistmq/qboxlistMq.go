package main

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"os"
	"path"
	"sync"

	log "qiniupkg.com/x/log.v7"

	"errors"
	"qbox.us/cc"
	"qbox.us/largefile"
	"strconv"
	"strings"

	"bytes"
	"fmt"
	"github.com/qiniu/xlog.v1"
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

//func write(msg string) {
//	fname := path.Join("/disk12/rsedit_save", "fh_error") //文件夹已经存在
//	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
//	if err != nil {
//		log.Fatalln("os.OpenFile: ", err)
//	}
//	_, err = f.WriteString(msg+"\n")
//	if err != nil {
//		log.Fatalln("f.WriteString", err)
//	}
//	f.Close()
//}

func filterMsg2(wg *sync.WaitGroup, in, out chan []byte) {
	defer wg.Done()
	xl := xlog.NewDummy()

	fname := path.Join("/disk12/rsedit_save", "fh_error") //文件夹已经存在
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
	if err != nil {
		log.Fatalln("os.OpenFile: ", err)
	}
	defer f.Close()

	for msg := range in {

		//读取 parse rsedit的消息
		m, err := ParseMessage("", string(msg))
		if err != nil {
			log.Fatal(err)
		}

		////bddelete过滤
		////op过滤
		//if !hasString(deleteFileOps, m.Op) {
		//	fmt.Println("m.op", m.Op)
		//	continue
		//}
		////bddelete消息内容
		//deleteMsg := Delete(m)
		//out <- []byte(deleteMsg)

		if m.Op == "delete" && len(m.Fh) == 0 {
			xl.Warn("ParseMessage len(m.Fh)==0 op delete")
			_, err = f.WriteString(string(msg) + "\n")
			if err != nil {
				log.Fatalln("f.WriteString", err)
			}
		}
		out <- msg
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
	endOffset := flag.Int64("e", 0, "end offset")
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
	tmpLen := int64(1024 * 512)
	fmt.Println("tmp", tmpLen)
	_, err = io.ReadFull(r, tmp)
	if err != nil {
		log.Fatal(err)
	}
	var realStartOffset int64
	var i int64
	bufr := bytes.NewReader(tmp)

	for ; i < tmpLen; i++ {
		fmt.Println("i: ", i)
		var h recordHeader
		bufr.Seek(i, 0)
		err = binary.Read(bufr, binary.LittleEndian, &h)
		if err != nil {
			fmt.Println("binary.Read, err :", err)
		}
		if h.MsgLen > uint16(tmpLen-i-8) {
			continue
		}
		if h.Tag != 0x0a {
			continue
		}
		//bufr
		m, err := ParseMessage("", string(tmp[i+8:i+8+int64(h.MsgLen)]))
		fmt.Println("err:", err)
		if err == nil {
			fmt.Printf("%+v\n", m)
			realStartOffset = *startOffset + i
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
			if err == io.EOF || r.Offset >= *endOffset {
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
