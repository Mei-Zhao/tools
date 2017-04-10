package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	//"os"
	"strconv"
	"strings"

	"github.com/qiniu/http/httputil.v1"
	"github.com/qiniu/log.v1"
	"github.com/qiniu/xlog.v1"
	bucketinfo2 "qbox.us/api/qconf/bucketinfo.v2"
	"qbox.us/api/qconf/domaing"
	"qbox.us/cc/config"
	"qbox.us/io/bucketinfo.v2"
	. "qbox.us/io/ucinfo"
	"qbox.us/memcache.v1"
	qconf "qbox.us/qconf/qconfapi"
)

func buildMemcacheKey(owner uint32, bucket, key string) string {

	return "io:" + strconv.FormatUint(uint64(owner), 36) + ":" + bucket + ":" + key
}

type memcacheValue struct {
	Fhandle  []byte `json:"f"`
	MimeType string `json:"m,omitempty"`
	AttName  string `json:"a,omitempty"`
	EndUser  string `json:"c,omitempty"`
	Fsize    int64  `json:"s"`
	KeyHint  uint32 `json:"k"`
	Uid      uint32 `json:"u"`
	PutTime  int64  `json:"t"`
}

//输入url,刷新缓存

type Config struct {
	DebugLevel int             `json:"debug_level"`
	McConns    []memcache.Conn `json:"mc_conns"`
	IoDomains  []string        `json:"io_domains"`
}

func main() {
	var conf Config
	config.Init("f", "qbox", "io_mc.conf")
	if err := config.Load(&conf); err != nil {
		log.Fatal("config.Load failed:", err)
	}

	log.SetOutputLevel(conf.DebugLevel)

	service := NewService(&conf)
	var url string
	_, err := fmt.Scan(&url)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("url: ", url)
	xl := xlog.NewDummy()
	cachekey := service.getCacheKey(xl, url)
	ret2 := memcacheValue{}
	err = service.Memcache.Get(xl, cachekey, &ret2)
	if err != nil {
		log.Println(cachekey, err)
		return
	}
	b, err := json.Marshal(ret2)
	if err != nil {
		log.Println("marshal json failed", err)
		return
	}
	fmt.Printf("%v %v\n", cachekey, string(b))
	var del string
	fmt.Print("\ndel?(y/n): ")
	_, err = fmt.Scan(&del)
	if err != nil {
		log.Println(err)
		return
	}
	if del != "y" && del != "Y" {
		return
	}
	err = service.Memcache.Del(xl, cachekey)
	if err == nil {
		fmt.Println("del success")
	} else {
		fmt.Println("del failed", err)
	}

}

type Service struct {
	Memcache memcache.Memcache
	bucketinfo.BucketGetter
	bucketinfo.BucketInfoGetter
	IoDomains []string
}

func NewService(conf *Config) *Service {

	mcS := memcache.New(conf.McConns)

	var qconfgS qconf.Config
	err := json.Unmarshal([]byte(qconfg), &qconfgS)
	if err != nil {
		log.Fatal(err)
	}
	qconfgcli := qconf.New(&qconfgS)

	bucketGetter := domaing.Client{qconfgcli}
	bucketInfoGetter := bucketinfo2.Client{qconfgcli}

	return &Service{
		Memcache:         mcS,
		BucketGetter:     bucketGetter,
		BucketInfoGetter: bucketInfoGetter,
		IoDomains:        conf.IoDomains,
	}
}

func UnEscape(in string) string {

	if len(in) < 2 {
		return in
	}
	nUnesc := 0
	var preCh = in[0]
	needUnEsc := preCh == '@'
	for i := 1; i < len(in); i++ {
		if needUnEsc && (in[i] == '/' || in[i] == '@') {
			nUnesc++
		}
		needUnEsc = (preCh == '/' && in[i] == '@')
		preCh = in[i]
	}
	if nUnesc == 0 {
		return in
	}

	var retBytes = make([]byte, 0, len(in)-nUnesc)

	preCh = in[0]
	needUnEsc = preCh == '@'
	for i := 1; i < len(in); i++ {
		if !(needUnEsc && (in[i] == '/' || in[i] == '@')) {
			retBytes = append(retBytes, preCh)
		}
		needUnEsc = (preCh == '/' && in[i] == '@')
		preCh = in[i]
	}
	retBytes = append(retBytes, preCh)

	return string(retBytes)
}

func (p *Service) getDomain(host string, rawurl string) (domain, key string) {

	host = strings.ToLower(host)

	for _, ioHost := range p.IoDomains {
		if host == ioHost {
			idx := strings.Index(rawurl, "/")
			if idx == -1 {
				return "/" + rawurl, ""
			}
			return "/" + rawurl[:idx], rawurl[idx+1:]
		}
	}
	return host, rawurl
}

func (service *Service) getCacheKey(xl *xlog.Logger, url string) string {
	index := strings.Index(url, "//")
	url = url[index+2:]
	params := strings.SplitN(url, "/", 2)
	host := params[0]
	realPath := params[1]
	xl.Info("domain: ", host, ",path: ", realPath)

	domain, relPath := service.getDomain(host, realPath)

	relPath = UnEscape(relPath)

	//通过域名查询空间信息必定是真实用户的信息
	domainInfo, err := service.BucketGetter.Get(xl, domain)
	var dotIndex int
	dotIndex = strings.Index(host, ".")
	if httputil.DetectCode(err) == 404 && dotIndex > 0 {
		domainInfo, err = service.BucketGetter.Get(xl, domain[dotIndex:])
	}
	if err != nil {
		xl.Fatalf("GetBucketInfo, domaing: Get failed, domain: %v, err: %v\n", domain, err)
	}

	bucketInfo, err := service.BucketInfoGetter.GetBucketInfo(xl, domainInfo.Uid, domainInfo.Tbl)
	if err != nil {
		xl.Fatalf("GetBucketInfo, buck_info: Get bucketInfo failed, uid: %v, tbl: %v, err: %v\n", domainInfo.Uid, domainInfo.Tbl, err)
	}

	bucketAllInfo := &BucketTblInfo{
		BucketInfo: bucketInfo,
		Info:       domainInfo,
	}

	ret := parseGetPath(relPath, bucketAllInfo)
	cachekey := buildMemcacheKey(bucketAllInfo.Uid, bucketAllInfo.Tbl, ret.key)

	if cachekey != "" {
		if bucketAllInfo.PreferStyleAsKey {
			xl.Println("PreferStyleAsKey:", bucketInfo.PreferStyleAsKey)
			cachekey = buildMemcacheKey(bucketAllInfo.Uid, bucketAllInfo.Tbl, relPath)
		}
	}

	xl.Println("uid: ", bucketAllInfo.Uid, ",bucket:", bucketAllInfo.Tbl, ",sourcekey:", ret.key, ",mc_key:", cachekey)
	return cachekey

}

type parseGetPathRet struct {
	key        string
	op         string // <fop>/<params>
	styleText  string // <sep><style>
	styleParam string
	access     bool
	addEndUser bool
	isKeyBase  bool // url里面只有keybase，没有完整的key
}

// Same as processOp in pub_svr.go
func parseGetPath(path string, info *BucketTblInfo) (ret parseGetPathRet) {

	ret.key = path
	ret.access = info.Protected != 1
	if info.Separator == "" {
		return
	}

	ret.addEndUser = info.Protected == 1

	n := strings.LastIndexAny(ret.key, info.Separator)
	if n < 0 {
		return
	}

	ret.styleText = ret.key[n:]
	style := ret.key[n+1:]
	if pos := strings.LastIndex(style, "@"); pos != -1 {
		ret.styleParam = style[pos+1:]
		style = style[:pos]
	}

	v, ok := info.Styles[style]
	if !ok { // no such style, skip
		ret.styleText = ""
		return
	}

	ret.key = ret.key[:n]
	ret.op = v
	ret.access = true

	if strings.HasPrefix(ret.op, "$0") {
		var keySuffix string
		if idx := strings.Index(ret.op, "?"); idx >= 0 {
			keySuffix, ret.op = ret.op[2:idx], ret.op[idx+1:]
		} else {
			keySuffix, ret.op = ret.op[2:], ""
		}
		if len(keySuffix) > 0 {
			ret.isKeyBase = true
			ret.key += keySuffix
		}
	}
	return
}

var qconfg = `{
    "mc_hosts": ["192.168.35.27:11215", "192.168.34.27:11215"],
    "master_hosts": ["http://192.168.34.29:8510", "http://192.168.34.30:8510"],
    "access_key": "ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ",
    "secret_key": "kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K",
    "lc_expires_ms": 300000,
    "lc_duration_ms": 5000,
    "lc_chan_bufsize": 16000,
    "mc_rw_timeout_ms": 100
    }
`
