package sanJiaoMaoTCPNet

import (
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"testing"
)

func TestSrv(t *testing.T) {
	srv := NewServer(defaultServerAdderss, defaultServerKey)

	srv.AddHandlerFunc("test", func(msg *Message.Message) {
		log.Println("server handler fun:test", msg.String())
		msg.SetResponseString(200, "respone success")
	})

	srv.Run()
}
