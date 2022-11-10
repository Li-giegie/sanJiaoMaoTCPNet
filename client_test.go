package sanJiaoMaoTCPNet

import (
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"sync"
	"testing"
	"time"
)

var c ClientI

var sy sync.WaitGroup

func TestClient1(t *testing.T) {

	n := 10000
	c = NewClient("127.0.0.1:9999", "client1")

	err := c.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	t1 := time.Now()
	for i := 0; i < n; i++ {
		sy.Add(1)
		go func(j int) {
			defer sy.Done()
			var testBuf = make([]byte, 1024, 1024)

			res, err := c.SendMessage("client2", "test", 200, testBuf, time.Second*2)
			if err != nil {
				t.Error("client sendMessage err:", err, res)
			}

			if j%n/10 == 0 {
				log.Println(res, string(res.Data))
			}
		}(i)
	}
	sy.Wait()
	fmt.Println(n, "次 ", "总耗时 ", time.Since(t1))
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
