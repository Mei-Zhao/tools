package main

import (
	"encoding/json"
	//"flag"
	//"fmt"
	//"os"
	"runtime"
	//"strings"
	//"sync"
	//"sync/atomic"
	//
	"github.com/qiniu/log.v1"
	//"github.com/qiniu/xlog.v1"
	//"github.com/tecbot/gorocksdb"
	//bson2 "gopkg.in/mgo.v2/bson"
	ebdtypes "qbox.us/ebd/api/types"
	"qbox.us/fh/fhver"
	"qbox.us/pfd/api/types"
	"qbox.us/pfdtracker/stater"
	"qbox.us/qconf/qconfapi"
	//"qbox.us/rocksdb/api/bsonmgr"
	//"qbox.us/rocksdb/api/rocksdb"
	"flag"
	"qbox.us/utils"
	"sync"
	"sync/atomic"

	"github.com/qiniu/xlog.v1"
	"qbox.us/rocksdb/api/client"
)

//nb tracker
var confJson = `{
   "mc_hosts": ["10.30.21.41:1141", "10.30.21.42:1141"],
      "master_hosts":  ["http://10.30.21.23:9020","http://10.30.21.24:9020"],
      "lc_expires_ms": 300000,
      "lc_duration_ms": 60000,
      "lc_chan_bufsize": 16000,
      "mc_rw_timeout_ms": 100
}`

func main() {
	var port string
	flag.StringVar(&port, "p", "", "rocksdb port")

	////param
	//var thread, ver int
	//var dbpath, bsons string
	//var writeLocal bool
	//
	//flag.StringVar(&dbpath, "d", "", "offset dbpath")
	//flag.StringVar(&bsons, "b", "", "bsondirs")
	//flag.BoolVar(&writeLocal, "w", false, "need writeLocal")
	//flag.IntVar(&ver, "v", 32, "offset db version")
	//flag.IntVar(&thread, "t", 50000, "size of channel")
	flag.Parse()
	runtime.GOMAXPROCS(24)
	host := "http://127.0.0.1:" + port
	cli := client.NewClient(host)

	var resultBuf chan client.ListFhRet
	//var errPatialBuf chan error
	xl := xlog.NewDummy()
	resultBuf, _ = cli.ListFhs(5000, "", xl)

	//
	//var conc int
	//if thread == 0 {
	//	conc = 50000
	//
	//} else {
	//	conc = thread
	//}
	dgids := [...] uint32 {2147483648,2147483649,2147483650,2147483651,2147483652,2147483653,2147483654,2147483655,2147483656,2147483657,2147483658,2147483692,2147483693,2147483694,2147483695,2147483696,2147483697,2147483698,2147483699,2147483700,2147483701,2147483702,2147483714,2147483715,2147483716,2147483717,2147483718,2147483719,2147483720,2147483721,2147483722,2147483723,2147483724,2147483736,2147483737,2147483738,2147483739,2147483740,2147483741,2147483742,2147483743,2147483744,2147483745,2147483746,2147483703,2147483704,2147483705,2147483706,2147483707,2147483708,2147483709,2147483710,2147483711,2147483712,2147483713,2147483747,2147483748,2147483749,2147483750,2147483751,2147483752,2147483753,2147483754,2147483755,2147483756,2147483757,2147483659,2147483660,2147483661,2147483662,2147483663,2147483664,2147483665,2147483666,2147483667,2147483668,2147483669,2147483681,2147483682,2147483683,2147483684,2147483685,2147483686,2147483687,2147483688,2147483689,2147483690,2147483691,2147483725,2147483726,2147483727,2147483728,2147483729,2147483730,2147483731,2147483732,2147483733,2147483734,2147483735,2147483670,2147483671,2147483672,2147483673,2147483674,2147483675,2147483676,2147483677,2147483678,2147483679,2147483680}
	var dgidSlice []uint32 = dgids[:]
	//
	//fmt.Printf("dbpath: %s, bsondirs: %s, type: %s, writeLocal: %v, ver: %v, thread: %v\n", dbpath, bsons, writeLocal, ver, thread)
	//
	//var offset []byte
	//fdNew := GetFd(ver)
	//offset = BuildValue(fdNew, 0)
	//
	//log.Infof("offset", offset)
	//
	////bsonmgr
	//var bsonDir []string
	//bsonDir = append(bsonDir, bsons)
	//bmgr, err := bsonmgr.NewEx(bsonDir)
	//if err != nil {
	//	log.Fatal(" bsonmgr.NewEx", bsonDir)
	//}
	//
	////tracker
	var conf qconfapi.Config
	err := json.Unmarshal([]byte(confJson), &conf)
	if err != nil {
		log.Fatal("json.Unmarshal", err)
	}
	gidStater := stater.NewGidStater(&conf)
	//
	//base := "/home/qboxserver/filter"
	//err = os.Mkdir(base, 0755)
	//if err != nil && os.IsNotExist(err) {
	//	log.Fatal("os.Mkdir", err)
	//}
	//
	////ebd
	//ebdBsonFile, err := utils.GetWriter("/home/qboxserver/filter", "ebd", os.O_WRONLY|os.O_APPEND)
	//if err != nil {
	//	log.Fatal("utils.GetWriter", err)
	//}
	//var ebdBsonFileLock sync.RWMutex
	//
	////open offsetDB withReadonly mode
	//ropt := gorocksdb.NewDefaultOptions()
	//ropt.SetWriteBufferSize(128 * 1024 * 1024)
	//rdb, err := gorocksdb.OpenDbForReadOnly(ropt, dbpath, true)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//resultBuf := putCh(offset, rdb, bmgr, conc)
	//
	var fileCount, total int64
	//
	var wg sync.WaitGroup
	readCh := func() {
		defer wg.Done()
		for {
			ret, ok := <-resultBuf
			if !ok {
				return
			}
			newTotal := atomic.AddInt64(&total, 1)
			if newTotal%100000 == 0 {
				log.Print("total:", newTotal)
			}
			//fh版本
			fhType := fhver.FhVer(ret.Fh)
			if fhType == 1 || fhType == 2 || fhType == 7 {
				continue
			}
			fhi, err2 := ebdtypes.DecodeFh(ret.Fh)
			if err2 != nil {
				log.Fatal("ebdtypes.DecodeFh", err2)
			}
			egid := types.EncodeGid(fhi.Gid)

			//一个文件一个reqid
			filexl := xlog.NewDummy()
			dgid, isECed, err3 := gidStater.State(filexl, egid)
			if err3 != nil {
				filexl.Fatal("self.gidStater.State ", err3)
			}
			filexl.Debugf("isECed :%v, egid :%v, dgid :%v\n", isECed, egid, dgid)

			if isECed {
				continue
			}
			if utils.InUids(dgid, dgidSlice) {
				newFileCount := atomic.AddInt64(&fileCount, 1)
				if newFileCount%100000 == 0 {
					log.Print("rpfd:", newFileCount)
				}
			}
		}
	}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go readCh()
	}
	wg.Wait()
	log.Println("rpfd total", fileCount)
	log.Println("total total", total)
	return
}

type ListFhRet struct {
	Fh      []byte "fh"
	Key     string
	Itbl    string
	Size    int64  "fsize"
	Hash    string "hash"
	PutTime int64  "putTime"
	ReqId   string
	Id      string "_id"
	Offset  int64
}

////seek according to offset
//func putCh(offset []byte, rdb *gorocksdb.DB, bmgr *bsonmgr.BsonFileMgrs, threadSize int) (resultBuf chan ListFhRet) {
//	resultBuf = make(chan ListFhRet, threadSize)
//	SwapValue(offset)
//	ro := gorocksdb.NewDefaultReadOptions()
//	ro.SetFillCache(false)
//	i := int64(0)
//	xl := xlog.NewDummy()
//	log.Println("begin")
//	iter := rdb.NewIterator(ro)
//	var ret ListFhRet
//	f := func() {
//		for iter.Seek(offset); iter.Valid(); iter.Next() {
//			value2 := rocksdb.GetSliceDataAndFree(iter.Key())
//			if i == 0 {
//				log.Println("first real offset in offsetdb", offset)
//			}
//
//			SwapValue(value2)
//			_, _, value, err := bmgr.ReadOneEx(value2, xl)
//			if err != nil {
//				xl.Error("bmgr.ReadOneEx", err)
//				return
//			}
//
//			err = bson2.Unmarshal(value, &ret)
//			if err != nil {
//				xl.Fatal("bson2.Unmarshal", err)
//			}
//			itblkey := strings.SplitN(ret.Id, ":", 2)
//			key := itblkey[1]
//			itbl := itblkey[0]
//			ret.Itbl = itbl
//			ret.Key = key
//			if i == 0 {
//				log.Println("first realy marker", ret.Id)
//			}
//			resultBuf <- ret
//			i++
//			if i%10000 == 0 {
//				log.Println("progress", i)
//			}
//		}
//		log.Println("total", i)
//		close(resultBuf)
//	}
//	go f()
//	return
//}
//
//func SwapValue(value []byte) {
//	value[3], value[4], value[6], value[7] = value[7], value[6], value[4], value[3]
//}
//
//func GetFd(fd int) (fdNew int) {
//	fdNew = fd*256 + 0 //伪造rocksdb-fd最开始offset
//	return
//}
//
//func BuildValue(fd int, offset int64) []byte {
//	value := make([]byte, 8)
//	value[0] = byte(fd)
//	value[1] = byte(fd >> 8)
//	value[2] = byte(fd >> 16)
//	value[3] = byte(offset)
//	value[4] = byte(offset >> 8)
//	value[5] = byte(offset >> 16)
//	value[6] = byte(offset >> 24)
//	value[7] = byte(offset >> 32)
//
//	return value
//}
