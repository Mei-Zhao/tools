//package main
//
//type OpLoG struct {
//	TS        int64  `json:"ts"`
//	H         int64  `json:"h"`
//	V         int    `json:"v"`
//	OP        string `json:"op"`
//	NS        string `json:"ns"`
//	RsInfo    RsInfo `json:"o"`
//	RsInfoNew RsInfo `json:"o2"`
//}
//
//type RsInfo struct {
//	Update   *RsInfo `json:"$set"` //set update
//	Id       string  `json:"_id"`  //insert or overwrite update
//	Fdel     int     `json:"fdel"`
//	Fh       []byte  `json:"fh"`
//	Key      string
//	Itbl     string
//	Size     int64  `json:"fsize"`
//	Hash     string `json:"hash"`
//	PutTime  int64  `json:"putTime"`
//	MimeType string `json:"mimeType"`
//	ReqId    string
//}
//
//func main() {
//	rsInfo := RsInfo{
//		Id:"7xtuph:hls4.v.momocdn.com/live/m_a007ebfda35279bc1479811592090107/online.m3u8",
//		Fh:
//	}
//	OpLoG := OpLoG{
//		OP: "u",
//		NS: "rs26.rs2",
//		V: 2,
//		TS: int64(6355753245264052538),
//		H: int64(794050438679424990),
//		RsInfo:rsInfo,
//	}
//
//}
//
////func TestDecodeOpLog(t *testing.T) {
////	oplog := `{"ts":6355753245264052538,"h":794050438679424990,"v":2,"op":"u","ns":"rs26.rs2","o":{"_id":"7xtuph:hls4.v.momocdn.com/live/m_a007ebfda35279bc1479811592090107/online.m3u8","fh":"Bpb/f/X0AABb37vqalqJFLmAh4IAAAAAyAMAAAAAAABGH+KBSQAAAOW9cHSnGFz8CfKZA/pg3bHMuGJb","fsize":968,"hash":"FuW9cHSnGFz8CfKZA_pg3bHMuGJb","ip":"AAAAAAAAAAAAAP//t4MHmQ==","mimeType":"application/vnd.apple.mpegurl","putTime":14798141194713003},"o2":{"_id":"7xtuph:hls4.v.momocdn.com/live/m_a007ebfda35279bc1479811592090107/online.m3u8"}}`
////	rsInfo := DecodeOpLog(oplog)
////	assert.Equal(t, "7xtuph", rsInfo.Itbl)
////	assert.Equal(t, "hls4.v.momocdn.com/live/m_a007ebfda35279bc1479811592090107/online.m3u8", rsInfo.Key)
////	assert.Equal(t, 0, rsInfo.Fdel)
////	fmt.Println(rsInfo.Fh)
////	bytes, _ := base64.StdEncoding.DecodeString("Bpb/f/X0AABb37vqalqJFLmAh4IAAAAAyAMAAAAAAABGH+KBSQAAAOW9cHSnGFz8CfKZA/pg3bHMuGJb")
////	fmt.Println("base64 decode []bytes", bytes)
////	fmt.Println("[]byte", []byte("Bpb/f/X0AABb37vqalqJFLmAh4IAAAAAyAMAAAAAAABGH+KBSQAAAOW9cHSnGFz8CfKZA/pg3bHMuGJb"))
////
////	//oplog = `{"ts":6351651530021535874,"h":9128004793423950324,"v":2,"op":"u","ns":"rs26.rs2","o":{"$set":{"mimeType":"image/bmp"}},"o2":{"_id":"7xlvdy:eMQCa4UqFYFB0Rac"}}`
////	//rsInfo = DecodeOpLog(oplog)
////	//assert.Equal(t, "7xlvdy", rsInfo.Itbl)
////	//assert.Equal(t, "eMQCa4UqFYFB0Rac", rsInfo.Key)
////	//assert.Equal(t, 0, rsInfo.Fdel)
////	//fmt.Println(rsInfo.Fh)
////	//
////	//oplog = `{"ts":6356057715495666028,"h":1217559487469680141,"v":2,"op":"i","ns":"rs26.rs2","o":{"_id":"77g9ya:7/583366560_1_65ab2104-bf9b-4018-8c3c-40a7725c196f_6","fh":"Bpb/fwEiAABTGaWkVZWJFLX5mi0IAAAA8FAJAAAAAAAeBU+rSQAAALK18yPHeQiuU07rE8SdNsN4vIjy","fsize":610544,"hash":"FrK18yPHeQiuU07rE8SdNsN4vIjy","ip":"AAAAAAAAAAAAAP//t4iNYg==","mimeType":"application/octet-stream","putTime":14798850093016294}}`
////	//rsInfo = DecodeOpLog(oplog)
////	//assert.Equal(t, "77g9ya", rsInfo.Itbl)
////	//assert.Equal(t, "7/583366560_1_65ab2104-bf9b-4018-8c3c-40a7725c196f_6", rsInfo.Key)
////	//assert.Equal(t, 0, rsInfo.Fdel)
////	//fmt.Println(rsInfo.Fh)
////	//
////	//oplog = `{"ts":6356057139970048422,"h":2542208009451632791,"v":2,"op":"u","ns":"rs26.rs2","o":{"_id":"7xylb6:7/660381580_1_c373dcf5-27db-4c05-8c4f-dffd7b2d1310","fdel":1,"fh":"Bpb/f/8qAAB5IF3J4PuGFO2Z9EQAAAAAoIkEAAAAAACwwDEKSAAAAPyFJTjg6QCNG/81H4Ay9cT1/W+s","fsize":297376,"hash":"FvyFJTjg6QCNG_81H4Ay9cT1_W-s","ip":"AAAAAAAAAAAAAP//c+e0zg==","mimeType":"application/octet-stream","putTime":14798848754197098},"o2":{"_id":"7xylb6:7/660381580_1_c373dcf5-27db-4c05-8c4f-dffd7b2d1310"}}`
////	//rsInfo = DecodeOpLog(oplog)
////	//assert.Equal(t, "7xylb6", rsInfo.Itbl)
////	//assert.Equal(t, "7/660381580_1_c373dcf5-27db-4c05-8c4f-dffd7b2d1310", rsInfo.Key)
////	//assert.Equal(t, 1, rsInfo.Fdel)
////	//fmt.Println(rsInfo.Fh)
////
////}
