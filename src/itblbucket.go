package main

import (
	"strconv"
	"encoding/json"
	"os"
	"io/ioutil"

	"github.com/qiniu/xlog.v1"
	qrpc "github.com/qiniu/rpc.v1"
	"qbox.us/api/tblmgr"
)


func main () {
	tblMgrHost := "http://192.168.34.200:10220"
	filexl := xlog.NewDummy()
	itbl, err3 := strconv.ParseUint(os.Args[1],36, 64)
	if err3 != nil {
		filexl.Fatal("strconv.ParseUint", err3)
	}
	itbl2 := strconv.FormatUint(itbl, 10)
	filexl.Infof("itbl2:", itbl2)
	url := tblMgrHost + "/itblbucket/" + itbl2
	resp, err3 := qrpc.DefaultClient.Get(filexl, url)
	if err3 != nil {
		filexl.Fatal("qrpc.DefaultClient.Get ", err3)
	}

	//读取body,更新缓存
	var entry2 tblmgr.BucketEntry
	var body []byte
	if resp.StatusCode != 200 && resp.StatusCode != 612 {
		filexl.Fatal("qrpc.DefaultClient.Get tblmgr ", resp.StatusCode)
	}
	if resp.StatusCode == 200 {
		body, err3 = ioutil.ReadAll(resp.Body)
		if err3 != nil {
			filexl.Fatal("ioutil.ReadAll", err3)
		}
		resp.Body.Close()

		err3 = json.Unmarshal(body, &entry2)
		if err3 != nil {
			filexl.Fatal("json.Unmarshal", err3)
		}
	}

	//612 no such entry 跳过下载
	if resp.StatusCode == 612 {
		resp.Body.Close()
		filexl.Warn("no such entry")
	}
	filexl.Println("uid", entry2.Uid)
	filexl.Println("bucket", entry2.Tbl)

}
