package main

import (
	"log"
	"os"

	"encoding/json"
	"fmt"
	"github.com/qiniu/xlog.v1"
	"qbox.us/api/kmqcli"
	"time"
)


type OpLoG struct {
	TS int64     `json:"ts"`
	H  int64     `json:"h"`
	V  int       `json:"v"`
	OP string    `json:"op"`
	NS string    `json:"ns"`
	FhRet  ListFhRet `json:"o"`
}

type ListFhRet struct {
	Id      string `json:"_id"`
	Fdel    int    `json:"fdel"`
	Fh      []byte `json:"fh"`
	Key     string
	Itbl    string
	Size    int64  `json:"fsize"`
	Hash    string `json:"hash"`
	PutTime int64  `json:"putTime"`
	ReqId   string
}

type stringSlice []string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var name string
	if len(os.Args) == 2 {
		name = os.Args[1]
	}
	host := "http://10.200.20.26:14532"
	hosts := make([]string, 1)
	hosts = append(hosts, host)

	cfg := &kmqcli.Config{
		//AccessKey: "ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ",
		//SecretKey: "kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K",
		AccessKey:"4_odedBxmrAHiu4Y0Qp0HPG0NANCf6VAsAjWL_k9",
		SecretKey:"SrRuUVfDX6drVRvpyN8mv8Vcm9XnMZzlbDfvVfMe",
		Hosts:     hosts,
		TryTimes:  2,
	}
	cli := kmqcli.New(cfg)
	xl := xlog.NewDummy()
	uid := uint32(260637563)
	code ,queues, err := cli.GetQueuesByUid(uid, xl)
	fmt.Println("code",code)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	for _, queue := range queues {
		fmt.Println("name", queue.Name)
	}
	code, msgs, position, err := cli.ConsumeMessages(name, "@", 1, xl)
	fmt.Println("code", code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("position", position)
	var oplog OpLoG
	if len(msgs) <= 0 {
		fmt.Println("len ",len(msgs))
		return
	}
	err = json.Unmarshal([]byte(msgs[0]), &oplog)
	if err != nil {
		xl.Error("json.Unmarshal")
	}
	fmt.Println("msgs", msgs[0])
	fmt.Println("oplog", oplog)
	fmt.Println("op", oplog.OP)
	fmt.Println("id", oplog.FhRet.Id)
	fmt.Println("fdel", oplog.FhRet.Fdel)
	fmt.Println("size", oplog.FhRet.Size)
	fmt.Println("fh", oplog.FhRet.Fh)
	fmt.Println("putTime", oplog.FhRet.PutTime)
	fmt.Println("now", time.Now().Unix())
	fmt.Println("seconds", oplog.FhRet.PutTime/1e7)
	fmt.Println("delay", time.Since(time.Unix(oplog.FhRet.PutTime/1e7,0)))
	if ( time.Since(time.Unix(oplog.FhRet.PutTime/1e7,0)) > (time.Duration(2*3600)*time.Second) ) {
		fmt.Println("delay 2 hours")
	}
}
