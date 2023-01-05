package sanJiaoMaoTCPNet

import (
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var c ClientI

var sy sync.WaitGroup

func TestClient1(t *testing.T) {

	n := 60000
	c = NewClient("127.0.0.1:9999", "client1")

	err := c.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	var testBuf = make([]byte, 10, 10)
	var ccc int
	t1 := time.Now()
	for i := 0; i < n; i++ {
		sy.Add(1)
		go func(j int) {
			defer sy.Done()

			res, err := c.SendMessage("client2", "test", 200, testBuf, time.Second*5)
			if err != nil {
				fmt.Println("error:",j,err,res)
				ccc++
				os.Exit(1)
				return
			}
			//if j%n/10 == 0 {
			//	log.Println(j,res,string(res.Data))
			//}
		}(i)
	}
	sy.Wait()
	fmt.Println(ccc,n, "次 ", "总耗时 ", time.Since(t1))
}

func TestClient2(t *testing.T) {
	c := NewClient("127.0.0.1:9999", "client2")

	c.AddHandlerFunc("test", func(msg *Message.Message, res Message.ReplyMessageI) {
		//fmt.Println(string(msg.Data))
		res.String(200, "client2 收到" + strconv.Itoa(int(msg.SrcApi)))
	})
	err := c.Connect()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("connect success----------")
	c.Run()
}
