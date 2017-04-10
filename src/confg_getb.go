package main

import (
	"os"

	qconf "qbox.us/qconf/qconfapi"
	"encoding/json"
	"github.com/qiniu/xlog.v1"
	"qbox.us/api/one/domain"
	"qbox.us/api/qconf/domaing"
	//"qbox.us/cc/config"
)

var confJson = `{
		"mc_hosts": ["10.34.37.21:1141", "10.34.37.22:1141"],
		"master": {
			"default": {
				"hosts": ["http://192.168.34.29:8510", "http://192.168.34.30:8510"],
				"fail_retry_interval_s": 10,
				"transport": {
				    "dial_timeout_ms": 200
				}
			},
			"failover": {
				"hosts": ["http://192.168.34.29:8510", "http://192.168.34.30:8510"],
				"transport": {
					"dial_timeout_ms": 200,
					"proxys": ["http://10.34.33.21:9060","http://10.34.33.22:9060","http://10.34.33.23:9060","http://10.34.33.24:9060"]
				}
			}
		},
		"access_key": "ppca4hFBYQ_ykozmLUcSIJi8eLnYhFahE0OF5MoZ",
		"secret_key": "kc6oDxKD3TYoRq3lUoS41-e4qtNYWzBSQZmigm7K",
		"lc_expires_ms": 300000,
		"lc_duration_ms": 60000,
		"lc_chan_bufsize": 16000,
		"mc_rw_timeout_ms": 100

}
`

type Info struct {
	Phy              string `json:"phy" bson:"phy"`
	Tbl              string `json:"tbl" bson:"tbl"`
	Uid              uint32 `json:"uid" bson:"uid"`
	Itbl             uint32 `json:"itbl" bson:"itbl"`
	Refresh          bool   `json:"refresh" bson:"refresh"`
	Global           bool   `json:"global" bson:"global"`
	Domain           string `json:"domain" bson:"domain"`
	domain.AntiLeech `json:"antileech,omitempty" bson:"antileech,omitempty"`
}

type ScanConfig struct {
	Qconfg                  qconf.Config      `json:"qconfg"`
}

func main() {
	id := os.Args[1]
	xl := xlog.NewDummy()
	xl.Info("domain: ", id)

	//config.Init("f", "qbox", "qboxscanbd.conf")
	//var conf ScanConfig
	//if err := config.Load(&conf); err != nil {
	//	xl.Fatal("config.Load failed:", err)
	//}
	//
	//s, err := NewScanAgent(&conf)
	//if err != nil {
	//	xl.Error("err", err)
	//}
	//ret, err := s.bucketGetter.Get(xl, id)
	//if err != nil {
	//	xl.Error("err", err)
	//}
	//xl.Printf("ret: %v\n", ret)

	var conf qconf.Config
	err := json.Unmarshal([]byte(confJson), &conf)
	if err != nil {
		xl.Fatal(err)
	}
	xl.Printf("conf: #v\n", conf)
	var ret Info
	qconfgcli := qconf.New(&conf)
	err = qconfgcli.Get(xl, &ret,id,1)
	if err != nil {
		xl.Error("qconfgcli.Get", err)
	}
	xl.Infof("ret: %#v\n", ret)

}

func NewScanAgent(cfg *ScanConfig) (*ScanAgent, error) {
	qconfgcli := qconf.New(&cfg.Qconfg)
	bucketGetter := domaing.Client{qconfgcli}
	return &ScanAgent{
		ScanConfig:    *cfg,
		bucketGetter: bucketGetter,
	}, nil
}

type ScanAgent struct {
	ScanConfig
	bucketGetter domaing.Client
}
