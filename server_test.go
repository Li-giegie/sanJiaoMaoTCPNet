package sanJiaoMaoTCPNet

import (
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"testing"
)

func TestSrv(t *testing.T) {

	srv := NewServer(defaultServerAdderss, defaultServerKey)
	//srv.AddAuthenticationHandle(func(key, ip string, text []byte) error {
	//	fmt.Println(key, ip, text)
	//	return errors.New("拒绝")
	//})
	srv.AddHandlerFunc("test", func(msg *Message.Message) {
		log.Println("server handler fun:test", msg.String())
		msg.SetResponseString(200, "respone success")
	})

	srv.Run()
}
