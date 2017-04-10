package main

import (
	"errors"
	"strconv"
	"strings"

	"qbox.us/cc/table"
	"qbox.us/pfd/api/types"
)

//const (
//	FhPfdLen = types.FileHandleBytes
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

// func (m *deleteMessage) String() string {
// return fmt.Sprintf("%s %d %d %s %s %d %d", base64.URLEncoding.EncodeToString(m.fh),
// m.ver, m.uid, m.bucket, m.key, m.modifyTime, m.putTime)
// }
//
//func (m *deleteMessage) serialize() string {
//	fields := []string{
//		string(m.fh), // fh
//		strconv.FormatUint(uint64(m.uid), 10),
//		table.Escape(m.bucket),
//		table.Escape(m.key),
//		strconv.FormatInt(m.modifyTime, 10),
//		strconv.FormatInt(m.putTime, 10),
//	}
//	return strings.Join(fields, "\t")
//}
//
//func (m *deleteMessage) current() bool {
//	return m.ver == 1
//}