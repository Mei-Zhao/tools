package main

import (
	"os"

	"qbox.us/api/uc"
	"qbox.us/mockacc"
)

const (
	TestUser   = "qboxtest"
	UcHost = "127.0.0.1:13200"
	USER_TYPE_STDUSER    = 0x0004
)
func main() {
	tr := mockacc.MakeTransport(mockacc.GetUid(TestUser), USER_TYPE_STDUSER)
	uc_cli := uc.New("http://"+UcHost, tr)

	sep := os.Args[1]
	fmt.p


}
