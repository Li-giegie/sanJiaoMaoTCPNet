package sanJiaoMaoTCPNet

import (
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"testing"
	"time"
)

var i = 1

func TestClient1(t *testing.T) {
	c := NewClient("127.0.0.1:9999", "client1")

	err := c.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	defer c.Close()

	var testBuf = make([]byte, 1024, 1024)

	for i := 0; i < 10000; i++ {
		res, err := c.SendMessage("client2", "test", 200, testBuf, time.Second*2)

		if err != nil {
			fmt.Println(i, err, string(res.Data))
			continue
		}
		if i != 9999 {
			continue
		}
		fmt.Println(res.String(), string(res.Data))
	}

	//c.Run()
	fmt.Println("Byte:1024B |request num:10000 |test success-------")
}

func TestClient2(t *testing.T) {
	c := NewClient("127.0.0.1:9999", "client2")

	c.AddHandlerFunc("test", func(msg *Message.Message, res Message.ReplyMessageI) {
		//fmt.Println(string(msg.Data))
		res.String(200, "client2 收到")
	})
	err := c.Connect()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("connect success----------")
	c.Run()
}
