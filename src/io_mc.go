package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/qiniu/xlog.v1"
	"qbox.us/api/one/domain"
	api3 "qbox.us/api/one/domain.v3"
	"qbox.us/memcache.v1"
	"qbox.us/api/qconf/bucketinfo.v2"
	qconf "qbox.us/qconf/qconfapi"
	. "qbox.us/io/ucinfo"
)

//{
//        "keys":["host_1_001", "host_1_002", "host_1_003"],
//        "host":"10.44.34.23:1121"
//},
//{
//        "keys":["host_2_001", "host_2_002", "host_2_003"],
//        "host":"10.44.34.30:1121"
//}//xs

//{
//"keys":["nb301_1_001", "nb301_1_002", "nb301_1_003"],
//"host":"192.168.35.27:11211"
//},
//{
//"keys":["nb301_2_001", "nb301_2_002", "nb301_2_003"],
//"host":"192.168.35.27:11212"
//},
//{
//"keys":["nb303_1_001", "nb303_1_002", "nb303_1_003"],
//"host":"192.168.35.29:11211"
//},
//{
//"keys":["nb303_2_001", "nb303_2_002", "nb303_2_003"],
//"host":"192.168.35.29:11212"
//},
//{
//"keys":["nb304_1_001", "nb304_1_002", "nb304_1_003"],
//"host":"192.168.35.30:11211"
//},
//{
//"keys":["nb304_2_001", "nb304_2_002", "nb304_2_003"],
//"host":"192.168.35.30:11212"
//},
//{
//"keys":["nb460_1_001", "nb460_1_002", "nb460_1_003"],
//"host":"192.168.35.42:11211"
//},
//{
//"keys":["nb460_2_001", "nb460_2_002", "nb460_2_003"],
//"host":"192.168.35.42:11212"
//},
//{
//"keys":["nb461_1_001", "nb461_1_002", "nb461_1_003"],
//"host":"192.168.35.43:11211"
//},
//{
//"keys":["nb461_2_001", "nb461_2_002", "nb461_2_003"],
//"host":"192.168.35.43:11212"
//},
//{
//"keys":["nb462_1_001", "nb462_1_002", "nb462_1_003"],
//"host":"192.168.35.44:11211"
//},
//{
//"keys":["nb462_2_001", "nb462_2_002", "nb462_2_003"],
//"host":"192.168.35.44:11212"
//},
//{
//"keys":["nb662_1_001", "nb662_1_002", "nb662_1_003"],
//"host":"192.168.44.55:11211"
//},
//{
//"keys":["nb662_2_001", "nb662_2_002", "nb662_2_003"],
//"host":"192.168.44.55:11212"
//},
//{
//"keys":["nb665_1_001", "nb665_1_002", "nb665_1_003"],
//"host":"192.168.44.58:11211"
//},
//{
//"keys":["nb665_2_001", "nb665_2_002", "nb665_2_003"],
//"host":"192.168.44.58:11212"
//},
//{
//"keys":["nb666_1_001", "nb666_1_002", "nb666_1_003"],
//"host":"192.168.44.59:11211"
//},
//{
//"keys":["nb666_2_001", "nb666_2_002", "nb666_2_003"],
//"host":"192.168.44.59:11212"
//},
//{
//"keys":["nb752_1_001", "nb752_1_002", "nb752_1_003"],
//"host":"192.168.34.40:11211"
//},
//{
//"keys":["nb752_2_001", "nb752_2_002", "nb752_2_003"],
//"host":"192.168.34.40:11212"
//},
//{
//"keys":["nb753_1_001", "nb753_1_002", "nb753_1_003"],
//"host":"192.168.34.46:11211"
//},
//{
//"keys":["nb753_2_001", "nb753_2_002", "nb753_2_003"],
//"host":"192.168.34.46:11212"
//}//nb

var mcs = `[
                //{
        	 //     "keys":["host_1_001", "host_1_002", "host_1_003"],
       		//       "host":"10.44.34.23:1121"
		//},
		//{
    		//        "keys":["host_2_001", "host_2_002", "host_2_003"],
       		//	 "host":"10.44.34.30:1121"
		//},
		//{
			//"keys":["nb301_1_001", "nb301_1_002", "nb301_1_003"],
			//"host":"192.168.35.27:11211"
		//},
		//{
			//"keys":["nb301_2_001", "nb301_2_002", "nb301_2_003"],
			//"host":"192.168.35.27:11212"
		//},
		//{
			//"keys":["nb303_1_001", "nb303_1_002", "nb303_1_003"],
			//"host":"192.168.35.29:11211"
		//},
		//{
			//"keys":["nb303_2_001", "nb303_2_002", "nb303_2_003"],
			//"host":"192.168.35.29:11212"
		//},
		//{
			//"keys":["nb304_1_001", "nb304_1_002", "nb304_1_003"],
			//"host":"192.168.35.30:11211"
		//},
		//{
			//"keys":["nb304_2_001", "nb304_2_002", "nb304_2_003"],
			//"host":"192.168.35.30:11212"
		//},
		//{
			//"keys":["nb460_1_001", "nb460_1_002", "nb460_1_003"],
			//"host":"192.168.35.42:11211"
		//},
		//{
			//"keys":["nb460_2_001", "nb460_2_002", "nb460_2_003"],
			//"host":"192.168.35.42:11212"
		//},
		//{
			//"keys":["nb461_1_001", "nb461_1_002", "nb461_1_003"],
			//"host":"192.168.35.43:11211"
		//},
		//{
			//"keys":["nb461_2_001", "nb461_2_002", "nb461_2_003"],
			//"host":"192.168.35.43:11212"
		//},
		//{
			//"keys":["nb462_1_001", "nb462_1_002", "nb462_1_003"],
			//"host":"192.168.35.44:11211"
		//},
		//{
			//"keys":["nb462_2_001", "nb462_2_002", "nb462_2_003"],
			//"host":"192.168.35.44:11212"
		//},
		//{
			//"keys":["nb662_1_001", "nb662_1_002", "nb662_1_003"],
			//"host":"192.168.44.55:11211"
		//},
		//{
			//"keys":["nb662_2_001", "nb662_2_002", "nb662_2_003"],
			//"host":"192.168.44.55:11212"
		//},
		//{
			//"keys":["nb665_1_001", "nb665_1_002", "nb665_1_003"],
			//"host":"192.168.44.58:11211"
		//},
		//{
			//"keys":["nb665_2_001", "nb665_2_002", "nb665_2_003"],
			//"host":"192.168.44.58:11212"
		//},
		//{
			//"keys":["nb666_1_001", "nb666_1_002", "nb666_1_003"],
			//"host":"192.168.44.59:11211"
		//},
		//{
			//"keys":["nb666_2_001", "nb666_2_002", "nb666_2_003"],
			//"host":"192.168.44.59:11212"
		//},
		//{
			//"keys":["nb752_1_001", "nb752_1_002", "nb752_1_003"],
			//"host":"192.168.34.40:11211"
		//},
		//{
			//"keys":["nb752_2_001", "nb752_2_002", "nb752_2_003"],
			//"host":"192.168.34.40:11212"
		//},
		//{
			//"keys":["nb753_1_001", "nb753_1_002", "nb753_1_003"],
			//"host":"192.168.34.46:11211"
		//},
		//{
			//"keys":["nb753_2_001", "nb753_2_002", "nb753_2_003"],
			//"host":"192.168.34.46:11212"
		//},
		//{
			//"keys":["bc21_1_001", "bc21_1_002", "bc21_1_003"],
			//"host":"10.30.21.41:1121"
		//},
		//{
			//"keys":["bc21_2_001", "bc21_2_002", "bc21_2_003"],
			//"host":"10.30.21.41:1122"
		//},
		//{
			//"keys":["bc22_1_001", "bc22_1_002", "bc22_1_003"],
			//"host":"10.30.21.42:1121"
		//},
		//{
			//"keys":["bc22_2_001", "bc22_2_002", "bc22_2_003"],
			//"host":"10.30.21.42:1122"
		//},
		//{
			//"keys":["bc72_1_001", "bc72_1_002", "bc72_1_003"],
			//"host":"10.30.23.43:1121"
		//},
		//{
			//"keys":["bc72_2_001", "bc72_2_002", "bc72_2_003"],
			//"host":"10.30.23.43:1122"
		//},
		//{
			//"keys":["bc73_1_001", "bc73_1_002", "bc73_1_003"],
			//"host":"10.30.23.44:1121"
		//},
		//{
			//"keys":["bc73_2_001", "bc73_2_002", "bc73_2_003"],
			//"host":"10.30.23.44:1122"
		//},
		 //{
                //        "keys":["host_1_001", "host_1_002", "host_1_003"],
                //        "host":"10.42.34.23:1121"
                //},
                //{
                //        "keys":["host_2_001", "host_2_002", "host_2_003"],
                //        "host":"10.42.34.30:1121"
                //},
                //{
                //        "keys":["host_1_001", "host_1_002", "host_1_003"],
                //        "host":"10.44.34.23:1121"
                //},
                //{
                //        "keys":["host_2_001", "host_2_002", "host_2_003"],
                //        "host":"10.44.34.30:1121"
                //},
                //{
                //        "keys":["host_1_001", "host_1_002", "host_1_003"],
                //        "host":"10.40.34.37:1121"
                //},
                //{
                //        "keys":["host_2_001", "host_2_002", "host_2_003"],
                //        "host":"10.40.34.51:1121"
                //}
                {
"keys":["nb301_1_001", "nb301_1_002", "nb301_1_003"],
"host":"192.168.35.27:11211"
},
{
"keys":["nb301_2_001", "nb301_2_002", "nb301_2_003"],
"host":"192.168.35.27:11212"
},
{
"keys":["nb303_1_001", "nb303_1_002", "nb303_1_003"],
"host":"192.168.35.29:11211"
},
{
"keys":["nb303_2_001", "nb303_2_002", "nb303_2_003"],
"host":"192.168.35.29:11212"
},
{
"keys":["nb304_1_001", "nb304_1_002", "nb304_1_003"],
"host":"192.168.35.30:11211"
},
{
"keys":["nb304_2_001", "nb304_2_002", "nb304_2_003"],
"host":"192.168.35.30:11212"
},
{
"keys":["nb460_1_001", "nb460_1_002", "nb460_1_003"],
"host":"192.168.35.42:11211"
},
{
"keys":["nb460_2_001", "nb460_2_002", "nb460_2_003"],
"host":"192.168.35.42:11212"
},
{
"keys":["nb461_1_001", "nb461_1_002", "nb461_1_003"],
"host":"192.168.35.43:11211"
},
{
"keys":["nb461_2_001", "nb461_2_002", "nb461_2_003"],
"host":"192.168.35.43:11212"
},
{
"keys":["nb462_1_001", "nb462_1_002", "nb462_1_003"],
"host":"192.168.35.44:11211"
},
{
"keys":["nb462_2_001", "nb462_2_002", "nb462_2_003"],
"host":"192.168.35.44:11212"
},
{
"keys":["nb662_1_001", "nb662_1_002", "nb662_1_003"],
"host":"192.168.44.55:11211"
},
{
"keys":["nb662_2_001", "nb662_2_002", "nb662_2_003"],
"host":"192.168.44.55:11212"
},
{
"keys":["nb665_1_001", "nb665_1_002", "nb665_1_003"],
"host":"192.168.44.58:11211"
},
{
"keys":["nb665_2_001", "nb665_2_002", "nb665_2_003"],
"host":"192.168.44.58:11212"
},
{
"keys":["nb666_1_001", "nb666_1_002", "nb666_1_003"],
"host":"192.168.44.59:11211"
},
{
"keys":["nb666_2_001", "nb666_2_002", "nb666_2_003"],
"host":"192.168.44.59:11212"
},
{
"keys":["nb752_1_001", "nb752_1_002", "nb752_1_003"],
"host":"192.168.34.40:11211"
},
{
"keys":["nb752_2_001", "nb752_2_002", "nb752_2_003"],
"host":"192.168.34.40:11212"
},
{
"keys":["nb753_1_001", "nb753_1_002", "nb753_1_003"],
"host":"192.168.34.46:11211"
},
{
"keys":["nb753_2_001", "nb753_2_002", "nb753_2_003"],
"host":"192.168.34.46:11212"

	]
`

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

func main() {
	xl := xlog.NewDummy()
	var conf []memcache.Conn
	err := json.Unmarshal([]byte(mcs), &conf)
	if err != nil {
		log.Fatal(err)
	}
	mcS := memcache.New(conf)
	url := os.Args[1]
	xl.Info("url:", url)
	index := strings.Index(url, "//")
	url = url[index+2:]
	xl.Info("url:", url)
	params := strings.SplitN(url, "/", 2)
	domain := params[0]
	key := params[1]
	xl.Info("domain:", domain, ",key:", key)
	//domainReal, relPath := getDomain(domain, key)

	oneHost := "http://192.168.34.29:23200"
	bucketInfo, err := getByDomain2(xl, domain, oneHost)
	if err != nil {
		xl.Fatal("err: ", err)
	}

	uid := bucketInfo.Uid
	bucket := bucketInfo.Tbl


	//bucketInfoGetter := getQconfg()
	//info, err:= bucketInfoGetter.GetBucketInfo(xl, uid, bucket)
	//if err != err {
	//	xl.Fatal("err:", err)
	//}
	//
	//
	parseUid := strconv.FormatUint(uint64(uid), 36)
	//ret := parseGetPath(params[1], info)
	//xl.Info("Key", ret.key)
	//cachekey := buildMemcacheKey(uid, bucketInfo.Tbl, ret.key)
	//if cachekey != "" {
	//	if info.PreferStyleAsKey {
	//		cachekey = buildMemcacheKey(uid, bucketInfo.Tbl, params[1])
	//	}
	//}

	mc_key := buildMemcacheKey(uid, bucket, key)
	cachekey := mc_key

	xl.Println("uid: ", uid, ",uid 36:", parseUid, ",bucket:", bucket, ",key:", key, ",mc_key:", cachekey)

	ret2 := memcacheValue{}
	err = mcS.Get(xl, cachekey, &ret2)
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
	err = mcS.Del(xl, cachekey)
	if err == nil {
		fmt.Println("del success")
	} else {
		fmt.Println("del failed", err)
	}
}

type GetByDomainRet struct {
	PhyTbl    string `json:"phy"`
	Tbl       string `json:"tbl"`
	Uid       uint32 `json:"owner"`
	Itbl      uint32 `json:"itbl"`
	Refresh   bool   `json:"refresh"`
	Global    bool   `json:"global"`
	Domain    string `json:"domain"`
	AntiLeech `json:"antileech,omitempty" bson:"antileech,omitempty"`
}

type AntiLeech struct {
	ReferWhiteList []string `json:"refer_wl,omitempty" bson:"refer_wl"`
	ReferBlackList []string `json:"refer_bl,omitempty" bson:"refer_bl"`
	ReferNoRefer   bool     `json:"no_refer" bson:"no_refer"`
	AntiLeechMode  int      `json:"anti_leech_mode" bson:"anti_leech_mode"` // 0:off,1:wl,2:bl
	AntiLeechUsed  bool     `json:"anti_leech_used" bson:"anti_leech_used"` // 表示是否设置过,只要设置过了就应该一直为true
}

func getByDomain2(xl *xlog.Logger, domainName string, oneHost string) (domain.GetByDomainRet, error) {
	c := api3.New(oneHost, nil)
	ret, err := c.GetByDomain(xl, domainName)
	return ret, err
}

var qconfg = `
    "mc_hosts": ["192.168.35.27:11215", "192.168.34.27:11215"],
    "master_hosts": ["http://192.168.34.29:8510", "http://192.168.34.30:8510"],
    "access_key": "ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ",
    "secret_key": "kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K",
    "lc_expires_ms": 300000,
    "lc_duration_ms": 5000,
    "lc_chan_bufsize": 16000,
    "mc_rw_timeout_ms": 100
`

func getQconfg() bucketinfo.Client{
	var conf qconf.Config
	err := json.Unmarshal([]byte(qconfg), &conf)
	if err != nil {
		log.Fatal(err)
	}


	qconfgcli := qconf.New(&conf)
	bucketInfoGetter := bucketinfo.Client{qconfgcli}
	return bucketInfoGetter
}

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


type parseGetPathRet struct {
	key        string
	op         string // <fop>/<params>
	styleText  string // <sep><style>
	styleParam string
	access     bool
	addEndUser bool
	isKeyBase  bool // url里面只有keybase，没有完整的key
}

//
//func getDomain(domain, path string) (){
//	host = strings.ToLower(domain)
//
//}