package sanJiaoMaoTCPNet

import (
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	srv, err := NewServer("127.0.0.1:9999", "server")

	if err != nil {
		log.Fatalln("err:", err)
	}
	srv.SetAuthentication(func(ip string, key string, data []byte) (bool, string) {
		return true, "测试拒绝认证"
	})
	srv.AddHandleFunc("test", func(msg *Message.Message, reply Message.ReplyMessageI) {
		fmt.Println("test----", msg.String(), string(msg.Data))
		reply.String(200, "server 收到")
	})

	log.Println("server run ---------")
	srv.Run()

}
