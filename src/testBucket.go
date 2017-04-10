package main

import (
	"os"

	"qiniupkg.com/api.v7/kodo"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/qiniu/xlog.v1"
	//"qbox.us/api/up"
	//"qiniu.com/auth/qboxmac.v1"
	//"qiniupkg.com/api.v7/kodocli"
	"bytes"
	//"encoding/base64"
	//"flag"
	//"fmt"
	"io"
	//"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	//"os"
	//"path/filepath"
	//"strconv"
	//"strings"
	"time"

	//"qbox.us/api/up"

	//"github.com/qiniu/api/auth/digest"

	//. "github.com/qiniu/api/conf"
	//"time"

	//"github.com/qiniu/api/rs"

	"io/ioutil"
	"strings"
	"compress/gzip"
	"qbox.us/api/up"
)

func main() {

	//fileName := "/Users/zhaomei/qbox/tools/src/file/a"
	newBucket()
	//file, err := os.Open(fileName)
	//xl := xlog.NewDummy()
	//if err != nil {
	//	xl.Fatal(err)
	//}
	//err = bucket.Put(nil, nil, "test_test2", file, int64(-1), nil)
	//xl.Error(err)
}

func newBucket() (bucket kodo.Bucket) {

	//QINIU_KODO_TEST = os.Getenv("QINIU_KODO_TEST")
	//if skipTest() {
	//	println("[INFO] QINIU_KODO_TEST: skipping to test qiniupkg.com/api.v7")
	//	return
	//}

	//ak := "Dp4_J7bDCAkqAutdsQHejKf1YlQxTVR0d8AFqRjP"
	//sk := "FD1bjSZY_PzzPO5TTmYJUgrKwVwmykGndK8bblBd"
	//if ak == "" || sk == "" {
	//	panic("require ACCESS_KEY & SECRET_KEY")
	//}
	//kodo.SetMac(ak, sk)
	//
	//bucketName := "zhaomei-admin"
	//domain := " "
	//if bucketName == "" || domain == "" {
	//	panic("require test env")
	//}
	//client := kodo.NewWithoutZone(nil)
	//
	//return client.Bucket(bucketName)


	//adminAk := "ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ"
	//adminSk := "kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K"
	//uid := uint32(1380758040)
	//bucketName := "beimei"

	//loacl
	adminAk := "4_odedBxmrAHiu4Y0Qp0HPG0NANCf6VAsAjWL_k9"
	adminSk := "SrRuUVfDX6drVRvpyN8mv8Vcm9XnMZzlbDfvVfMe"
	uid := uint32(260637563)
	bucketName := "test"
	key := "test"
	UP_HOST := "http://up-na0.qiniu.com"
	//token := getToken(mac, fmt.Sprintf("%v/0", uid), nil)

	fileName := "file/a"
	file, err := os.Open(fileName)
	xl := xlog.NewDummy()
	if err != nil {
		xl.Fatal(err)
	}
	pr, pw := io.Pipe()
	defer pr.Close()
	w := gzip.NewWriter(pw)
	go func() {
		_, _ = io.Copy(w, file)
		w.Close()
		pw.Close()
	}()

	//policy := up.AuthPolicy{
	//	Scope: bucketName + ":" + key,
	//}
	putPolicy := up.AuthPolicy{
		Scope:  bucketName + ":" + key,
		Deadline: time.Now().Unix() + 3600,
	}
	token := makeUpToken(uid, putPolicy, adminAk, adminSk)
	//uploader := kodocli.NewUploaderWithoutZone(nil)
	//err = uploader.Put(nil, nil, token, "admin4", pr, int64(-1), nil)
	//xl.Error("err", err)
	//config := kodo.Config{
	//	AccessKey: adminAk,
	//	SecretKey: adminSk,
	//	Transport: transSu,
	//}
	//client := kodo.NewWithoutZone(&config)
	//return client.Bucket(bucketName)
	//service := rs.New(transSu)
	req, err := rsupload(UP_HOST, pr, key, token)
	if err != nil {
		log.Fatal(err)
	}
	//req.Host = UP_HOST

	log.Printf("req : %#v", req)
	fmt.Println(req.Host)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	dumpResp("up", resp)

	return
}

func dumpResp(service string, resp *http.Response) {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 30), service, " resp", strings.Repeat("-", 30))
	defer resp.Body.Close()
	fmt.Println("Code:", resp.StatusCode)
	fmt.Println("body", resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read body failed, %v", err)
	}

	fmt.Println("Header:")
	for k, v := range resp.Header {
		fmt.Printf("\t%v: %v\n", k, v)
	}

	fmt.Println("Body:")
	fmt.Printf("%v\n", string(body))
}

func rsupload(host string, f io.Reader, filename string, uptoken string) (*http.Request, error) {
	// POST UpHost
	// Content-Type: multipart/form-data; boundary=<Boundary>
	// Body:
	// 1. key = <Key>
	// 2. crc32 = <CRC32>
	// 3. token = <UpToken>
	// 4. file = <FileData>
	fmt.Println("Path:", "/")
	extraParams := map[string]string{
		"token": uptoken,
		"key":   filename,
	}
	return newfileUploadRequest(host, extraParams, "file", f, filename)
}

func newfileUploadRequest(url string, params map[string]string, paramName string, file io.Reader, filename string) (req *http.Request, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filename)
	//part, err := writer.CreateFormFile(paramName, "")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	// finally create the request
	req, err = http.NewRequest("POST", url, body)
	fmt.Println("url", url)
	if err != nil {
		return
	}
	// set type & boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//log.Println("Content-Type", writer.FormDataContentType())
	return
}

func makeUpToken(uid uint32, putPolicy up.AuthPolicy, adminAk, adminSk string) string {
	//putPolicy.Deadline = time.Now().Add(time.Minute).Unix()
	//
	//b, _ := json.Marshal(putPolicy)
	//data := base64.URLEncoding.EncodeToString(b)
	//
	//suInfo := fmt.Sprintf("%v/0", uid)
	//hash := hmac.New(sha1.New, []byte(adminSk))
	//hash.Write([]byte(suInfo))
	//hash.Write([]byte(data))
	//sign := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	//
	//token := fmt.Sprintf("%s:%s:%s:%s", suInfo, adminAk, sign, data)
	//fmt.Println("token", token)
	//putPolicy.Expires = uint32(time.Now().Add(time.Minute).Unix())

	b, _ := json.Marshal(putPolicy)
	data := base64.URLEncoding.EncodeToString(b)

	suInfo := fmt.Sprintf(":%d/0:", uid)
	hash := hmac.New(sha1.New, []byte(adminSk))
	hash.Write([]byte(suInfo))
	hash.Write([]byte(data))
	sign := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	token := fmt.Sprintf("%s%s:%s:%s", suInfo, adminAk, sign, data)
	fmt.Println("token", token)
	return token
}