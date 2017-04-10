package main

import (
	"time"
	"log"
)
func main() {

	timeStr := "2017-02-12T08:27:34.418543Z"
	timeClf, date := ParseTime(timeStr)
	log.Println("timeClf", timeClf)
	log.Println("date", date)
}
//上传的文件名中的时间是CST时间,目前线上机器都是CST时间
func ParseTime(timeStr string) (timeCLF string, date string) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatal("time.LoadLocation: Asia/Shanghai, err: ", err)
	}
	time, err := time.ParseInLocation(time.RFC3339, timeStr, time.UTC)
	if err != nil {
		log.Fatal("time.ParseInLocation ", timeStr, err)
	}
	time = time.In(loc)
	timeCLF = time.Format("02/Jan/2006:15:04:05 +0800")
	date = time.Format("2006-01-02-15-04")
	return
}