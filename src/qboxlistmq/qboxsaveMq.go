package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"time"

	"qbox.us/qmq/qmqapi/v1/mq"

	//"github.com/ngaut/log"
	"github.com/qiniu/api/auth/digest"
	xlog "github.com/qiniu/xlog.v1"
	//"fmt"
)

func main() {
	mqHost := flag.String("mqhost", "http://192.168.58.62:14500", "mq host")//默认2199 14500 rsedit mq
	flag.Parse()
	trans := digest.NewTransport(&digest.Mac{
		AccessKey: "-s_NZdOCsmL9vRsY_TAhxnBuDLh4dgh-HqM8eRV_",
		SecretKey: []byte("3b2ghfzvvfP8laMGb9z5c_zulNTZ1pzbHa008wDZ"),
	}, nil)
	mqcli := mq.NewWithHost(trans, *mqHost)
	decoder := json.NewDecoder(os.Stdin)
	var msg []byte
	xl := xlog.NewDummy()
	for {
		err := decoder.Decode(&msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			xl.Fatal(err)
		}
		success := false
		for i := 0; i < 20; i++ {
			_, err = mqcli.Put(xl, "rsedit", msg)
			if err == nil {
				success = true
				//fmt.Println("true")
				break
			}
			xl.Warn(err)
			time.Sleep(time.Second)
		}
		if !success {
			xl.Fatal(err)
		}
	}
}