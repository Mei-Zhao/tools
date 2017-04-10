package main

import (
	//"bytes"
	"compress/gzip"
	"fmt"
	"time"
	"os"

	//"bufio"
	//"io/ioutil"
	"io"
	"github.com/qiniu/xlog.v1"
	. "github.com/qiniu/api/conf"
	"qbox.us/api/up"
	"qbox.us/api/v2/rs"
	"qiniu.com/auth/qboxmac.v1"
)

func main() {
	ak := "4_odedBxmrAHiu4Y0Qp0HPG0NANCf6VAsAjWL_k9"
	sk := "SrRuUVfDX6drVRvpyN8mv8Vcm9XnMZzlbDfvVfMe"
	uphost := "http://127.0.0.1:7777"
	key := "testGzip_13"
	bucket := "bucket"
	fileName := "/Users/zhaomei/qbox/tools/src/file/a"
	progress := int64(0)
	rangeSzie := int64(5)
	uid := uint32(260637563)
	xl := xlog.NewDummy()
	file, err := os.Open(fileName)
	if err != nil {
		xl.Fatal(err)
	}
	err = upload(xl, ak, sk, uphost, key, bucket, uid, file, progress, rangeSzie)
	xl.Error(err)

}

func upload(xl *xlog.Logger, adminAk, adminSk, uphost, key, bucket string, uid uint32, file *os.File, process, rangeSize int64) error {
	UP_HOST = uphost
	mac := &qboxmac.Mac{adminAk, []byte(adminSk)}
	transSu := qboxmac.NewAdminTransport(mac, fmt.Sprintf("%v/0", uid), nil)
	policy := up.AuthPolicy{
		Deadline: time.Now().Unix() + 3600,
	}

	xl.Infof("will upload file as key :%s, for user :%v, bucket: %s\n", key, uid, bucket)
	entry := bucket + ":" + key
	policy.Scope = entry
	service := rs.New(transSu)
	pr, pw := io.Pipe()
	defer pr.Close()
	w := gzip.NewWriter(pw)
	go func() {
		_, _ = io.Copy(w, file)
		w.Close()
		pw.Close()
	}()
	var err error
	//f2, err := os.Create("gziped.gz")
	//defer f2.Close()
	//
	//_, err = io.Copy(f2, pr)

	err = service.Put2WithStream(nil, false, key, pr, int64(-1), "", "", nil)
	if err != nil {
		xl.Errorf("upload failed, err :%s, entryURI: %s, uphost :%s\n", err, entry, UP_HOST)
		return err
	}
	//_, code, err := service.Put(entry, "", pr, int64(34), "", "", "")
	//if err != nil || code != 200 {
	//	xl.Errorf("upload failed, code: %v, err :%s, entryURI: %s, uphost :%s\n", code, err, entry, UP_HOST)
	//	return err
	//}

	return nil
}
