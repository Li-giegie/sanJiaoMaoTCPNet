package sanJiaoMaoTCPNet

import (
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"testing"
	"time"
)

func TestClient2(t *testing.T) {
	cli := NewClient(defaultServerAdderss, defaultClientKey+"2")

	err := cli.Connect()

	if err != nil {
		log.Fatalln(err)
	}

	cli.AddHandlerFunc("c2test1", func(res Responser, msg *Message.Message) {
		msg.SetResponseString(200, "client2 success")
		res.Response(msg)
	})

	for {

	}
}

func TestClient1(t *testing.T) {
	cli := NewClient(defaultServerAdderss, defaultClientKey+"1")

	err := cli.Connect()

	if err != nil {

		log.Fatalln(err)
	}
	for {
		msg := Message.NewMsg(defaultClientKey + "1")
		msg.SetRequestString(defaultClientKey+"2", "c2test1", "hello ?")
		res, err := cli.Request(msg)
		if err != nil {
			log.Fatalln("request err:", err)
		}
		log.Println(res.String())

		time.Sleep(time.Second * 2)
	}

}
