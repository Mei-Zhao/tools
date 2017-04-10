package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/api/auth/digest"
)

var (
	mac = &digest.Mac{"ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ", []byte("kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K")}
	cli = digest.NewClient(mac, nil)
)

type RSInfo struct {
	ID         string    `json:"_id,omitempty"`
	Uid        string    `json:"uid"`
	Bucket     string    `json:"bucket"`
	Key        string    `json:"key"`
	Fsize      int64     `json:"fsize"`
	Hash       string    `json:"hash"`
	Fdel       *int      `json:"fdel"`
	PutTime    int64     `json:"putTime"`
	PutTimeStr time.Time `json:"putTimeStr"`
	MimeType   string    `json:"mimeType"`
	Fh         string    `json:"fh"`
	FhURL      string    `json:"fhURL"`
	Osize      *int64    `json:"osize,omitempty"`
	Itbl       int       `json:"itbl,omitempty"`
	IP         string    `json:"ip,omitempty"`
}

type Entry struct {
	Entry RSInfo
	Error string
	Name  string
}

type Entrys []Entry

func inspect(uid, bucket, key string, itbl int) {
	path := "/admin/inspect/" + base64.URLEncoding.EncodeToString([]byte(bucket+":"+key)) + "/user/" + uid
	host := fmt.Sprintf("https://rs.qbox.me")
	resp, err := cli.Get(host + path)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	entryInfo := Entrys{}
	err = json.Unmarshal(data, &entryInfo)
	if err != nil {
		log.Println(err, string(data))
		return
	}
	for _, s := range entryInfo {
		s.Entry.Uid = uid
		s.Entry.Bucket = bucket
		s.Entry.Key = key
		s.Entry.Itbl = itbl
		s.Entry.PutTimeStr = time.Unix(s.Entry.PutTime/1e7, (s.Entry.PutTime%1e7)*100)
		s.Entry.FhURL = strings.Replace(s.Entry.Fh, "/", "_", -1)
		s.Entry.FhURL = strings.Replace(s.Entry.FhURL, "+", "-", -1)
		if s.Error == "none" {
			fmt.Printf("===> %s %v\n", s.Name, ToJson(s.Entry))
		} else {
			fmt.Printf("===> %s\t%s\n", s.Name, s.Error)
		}
	}
}

type BucketInfo struct {
	Tbl   string
	Owner int
	Itbl  int `json:"itbl"`
}

func getByDomain(domain string) (uid, tbl string, itbl int) {
	u := "https://api.qiniu.com/v6/admin/domain/getbydomain?domain=" + domain
	resp, err := cli.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	bucketInfo := BucketInfo{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &bucketInfo)
	if err != nil {
		log.Fatal(err)
	}
	uid = strconv.Itoa(bucketInfo.Owner)
	tbl = bucketInfo.Tbl
	itbl = bucketInfo.Itbl
	return
}

func viewHeader(rawurl string) {
	log.Printf("------------HEAD------------\n")
	resp, err := http.Head(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP/1.1", resp.Status)
	for k, v := range resp.Header {
		for _, vv := range v {
			fmt.Printf("%s: %s\n", k, vv)
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var uid, tbl, key string
	var itbl int
	if len(os.Args) == 4 {
		uid, tbl, key = os.Args[1], os.Args[2], os.Args[3]
	} else if len(os.Args) == 3 {
		domainZ0 := os.Args[1] + ".com1.z0.glb.clouddn.com"
		domainZ1 := os.Args[1] + ".com1.z1.glb.clouddn.com"
		uid, tbl, itbl = getByDomain(domainZ0)
		if uid == "0" {
			uid, tbl, itbl = getByDomain(domainZ1)
		}
		if uid == "0" {
			log.Fatal("get uid and bucket from itbl failed:", os.Args[1], uid, tbl)
		}
		key = os.Args[2]
	} else if len(os.Args) == 2 {
		rawurl := os.Args[1]
		if !strings.HasPrefix(rawurl, "http://") {
			rawurl = "http://" + rawurl
		}
		u, err := url.Parse(rawurl)
		if err != nil {
			log.Fatal(err)
		}
		uid, tbl, itbl = getByDomain(u.Host)
		if len(u.Path) == 0 {
			key = ""
		} else {
			key = u.Path[1:]
		}
		if uid == "0" {
			log.Fatal("get uid from domain failed: ", rawurl)
		}
	}
	uidI, err := strconv.Atoi(uid)
	if err != nil {
		log.Fatal("invalid uid", uid)
	}
	log.Println(uid, tbl, key)
	fmt.Printf("io_mc_key:\nio:%s:%s:%s\n\n", strconv.FormatInt(int64(uidI), 36), tbl, key)
	inspect(uid, tbl, key, itbl)
}

func ToJson(v interface{}) string {
	// b, err := json.Marshal(v)
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return fmt.Sprintf("{\"error\":%s}", err.Error())
	}
	return string(b)
}