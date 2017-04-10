//package main
//
//import (
//	"bufio"
//	"encoding/binary"
//	"encoding/json"
//	"flag"
//	"io"
//	"os"
//	"sync"
//
//	log "qiniupkg.com/x/log.v7"
//
//	"qbox.us/cc"
//	"qbox.us/ebd/api/types"
//	"qbox.us/largefile"
//
//	"errors"
//	"strconv"
//	"strings"
//
//	"qbox.us/cc/table"
//	pfdtypes "qbox.us/pfd/api/types"
//)
//
//type recordHeader struct {
//	MsgLen   uint16
//	Tag      uint8
//	State    uint8
//	LockTime uint32 // in s
//}
//
//const (
//	indexSize     = 32
//	headerSize    = 8
//	recordTag     = 0x0a
//	offsetOfState = 3
//)
//
//const (
//	stateNormal     = 0x01
//	stateProcessing = 0x02
//	stateDone       = 0x04
//
//	stateTravelPut     = stateNormal | stateProcessing | stateDone
//	stateTravelGet     = stateProcessing | stateDone
//	stateTravelTimeout = stateDone
//)
//
//const (
//	FhPfdLen = pfdtypes.FileHandleBytes
//)
//
//type deleteMessage struct {
//	ver        int
//	fh         []byte
//	uid        uint32
//	bucket     string
//	key        string
//	modifyTime int64
//	putTime    int64
//	op         string
//}
//
//func parseDeleteMessage(msg []byte) (m *deleteMessage, err error) {
//	if len(msg) < FhPfdLen {
//		err = errors.New("msg is too short, invalid fh")
//		return
//	}
//
//	if len(msg) == FhPfdLen {
//		// for compatibility when migration, fh can only be FhPfd or FhPfdV2
//		m = &deleteMessage{
//			ver: 0,
//			fh:  msg,
//		}
//		return
//	}
//
//	var info deleteMessage
//	info.ver = 1
//
//	// has modifyTime & putTime:   fh \t uid \t bucket \t key \t modifyTime \t putTime \t op
//	info.fh = msg[:FhPfdLen]
//	fields := strings.Split(string(msg[FhPfdLen:]), "\t")
//
//	if len(fields) < 6 {
//		err = errors.New("invalid message")
//		return
//	}
//
//	// fields[0] is empty
//	var uid uint64
//	uid, err = strconv.ParseUint(fields[1], 10, 32)
//	if err != nil {
//		return
//	}
//	info.uid = uint32(uid)
//
//	info.bucket, err = table.Unescape(fields[2])
//	if err != nil {
//		return
//	}
//	info.key, err = table.Unescape(fields[3])
//	if err != nil {
//		return
//	}
//
//	info.modifyTime, err = strconv.ParseInt(fields[4], 10, 64)
//	if err != nil {
//		return
//	}
//	info.putTime, err = strconv.ParseInt(fields[5], 10, 64)
//	if err != nil {
//		return
//	}
//	if len(fields) >= 7 {
//		info.op, err = table.Unescape(fields[6])
//		if err != nil {
//			return
//		}
//	}
//
//	m = &info
//	return
//}
//
//func filterMsg(wg *sync.WaitGroup, in, out chan []byte) {
//	defer wg.Done()
//	for msg := range in {
//		m, err := parseDeleteMessage(msg)
//		if err != nil {
//			log.Fatal(err)
//		}
//		if m.uid != 1380260887 {
//			out <- msg
//		}
//	}
//}
//
//func main() {
//	startOffset := flag.Int64("s", 0, "start offset")
//	limit := flag.Int64("c", -1, "read max count")
//	dir := flag.String("d", "", "dir")
//	chunkBit := flag.Uint("b", 29, "chunk bit")
//	save := flag.String("save", "save", "save file")
//	process := flag.Int("p", 5, "process to parse")
//	flag.Parse()
//	f, err := largefile.Open(*dir, *chunkBit)
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := &cc.Reader{f, *startOffset}
//
//	// find real start offset
//	var tmp = make([]byte, 4096)
//	_, err = io.ReadFull(r, tmp)
//	if err != nil {
//		log.Fatal(err)
//	}
//	var realStartOffset, i int64
//	for ; i < 4096; i++ {
//		_, err = types.DecodeFh(tmp[i : i+60])
//		if err == nil {
//			realStartOffset = *startOffset + i - 8
//			log.Println("real start offset", realStartOffset)
//			break
//		}
//	}
//	r = &cc.Reader{f, realStartOffset}
//	br := bufio.NewReaderSize(r, 4*1024*1024)
//	saveFile, err := os.Create(*save)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer saveFile.Close()
//	encoder := json.NewEncoder(saveFile)
//	var readCh = make(chan []byte, 1)
//	var saveCh = make(chan []byte)
//	var wg = &sync.WaitGroup{}
//	wg.Add(*process)
//	for i := 0; i < *process; i++ {
//		go filterMsg(wg, readCh, saveCh)
//	}
//	go func() {
//		for msg := range saveCh {
//			err = encoder.Encode(msg)
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//	}()
//	for count := int64(0); count != *limit; count++ {
//		var h recordHeader
//		err = binary.Read(br, binary.LittleEndian, &h)
//		if err != nil {
//			if err == io.EOF {
//				close(readCh)
//				break
//			}
//			log.Fatal(err)
//		}
//		msg := make([]byte, h.MsgLen)
//		n, err := io.ReadFull(br, msg)
//		if err != nil {
//			if err == io.EOF {
//				break
//			}
//			log.Fatal(n, err)
//		}
//		readCh <- msg
//	}
//	wg.Wait()
//	close(saveCh)
//}