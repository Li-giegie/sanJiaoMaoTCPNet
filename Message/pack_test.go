package Message

import (
	"bytes"
	"fmt"
	"testing"
)

var testPackBuf =make([]byte,500)

func BenchmarkPack(b *testing.B) {

	var msg = Message{
		Header:Header{
			SrcKey:        "srcKey",
			SrcApi:        123,
			DistKey:       "distKey",
			DistApi:       "distApi",
			MType:         1,
			StateCode:     200,
			DataLen:       0,
		},
		Data:          testPackBuf,
	}

	for i := 0; i < b.N; i++ {
		buf,err := Pack(&msg)
		if err != nil {b.Error(err,buf);return }
		//fmt.Println(string(buf))
		//fmt.Println(len(buf))
	}
}

//比proto 多13个字节 少字节情况下性能不分上下 大字节占优势 与proto差距不大
func BenchmarkPackV2(b *testing.B) {

	var msg = Message{
		Header:Header{
			SrcKey:        "srcKey",
			SrcApi:        123,
			DistKey:       "distKey",
			DistApi:       "distApi",
			MType:         1,
			StateCode:     200,
			DataLen:       0,
		},
		Data:          testPackBuf,
	}
	for i := 0; i < b.N; i++ {
		buf,err := PackV2(&msg)
		if err != nil {b.Error(err,buf);return }
		//fmt.Println(len(buf))
	}
}

func BenchmarkUnPack(b *testing.B) {
	buf,err := Pack(&Message{
		Header:Header{
			SrcKey:        "srcKey",
			SrcApi:        123,
			DistKey:       "distKey",
			DistApi:       "distApi",
			MType:         1,
			StateCode:     200,
			DataLen:       0,
		},
		Data:          testPackBuf,
	})
	if err != nil {b.Error(err);return}
	var r *bytes.Buffer
	var msg *Message
	for i := 0; i < b.N; i++ {
		r = bytes.NewBuffer(buf)
		msg,err = UnPack(r)
		if err != nil { b.Error(err,msg);return}

		//fmt.Println(msg.String(),string(msg.Data))
	}

	//fmt.Println(msg.String(),string(msg.Data))
}

func BenchmarkUnPackV2(b *testing.B) {

	buf,err := PackV2(&Message{
		Header:Header{
			SrcKey:        "srcKey",
			SrcApi:        123,
			DistKey:       "distKey",
			DistApi:       "distApi",
			MType:         1,
			StateCode:     200,
			DataLen:       0,
		},
		Data:          testPackBuf,
	})

	if err != nil {b.Error(err);return}
	var r *bytes.Buffer
	var msg *Message
	for i := 0; i < b.N; i++ {
		r = bytes.NewBuffer(buf)
		msg,err = UnPackV2(r)
		if err != nil { b.Error(err,msg);return}
	}
	//fmt.Println(string(msg.Data))
}
func BenchmarkWrite(b *testing.B) { //7986858               185.0 ns/op

	var m = make(map[int64]string)
	for i := 0; i < b.N; i++ {
		m[int64(i)] = "asd"
	}
	fmt.Println(len(m))
}

func BenchmarkCopy(b *testing.B) {
	slice1 := []byte{1, 2, 3, 4, 5}
	slice2 := make([]byte,5)
	copy(slice2, slice1) // 只会复制slice1的前3个元素到slice2中
	//copy(slice1, slice2) // 只会复制s
	fmt.Println(slice1,slice2)
	return
	for i := 0; i < b.N; i++ {

	}
}