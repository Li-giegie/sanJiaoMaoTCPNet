package sanJiaoMaoTCPNet

import (
	"errors"
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/utils"
	"io"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	Request_ClientType byte = 1
	Respone_ClientType byte = 2
)

type Responser interface {
	Response(mes *Message.Message)
}
type Client struct {
	ip   net.IP
	prot int

	srcKey string

	conn *net.TCPConn

	TimeOut time.Duration

	HandlerFunc map[string]func(res Responser, msg *Message.Message)

	responeMsg map[int64]chan *Message.Message

	cleanMemoryFragment int

	AuthenticationText []byte

	PushHanderFun func(msg *Message.Message)
}

func NewClient(remoteAdderss, srckey string) *Client {

	var cli = Client{
		srcKey:              srckey,
		TimeOut:             3,
		HandlerFunc:         map[string]func(res Responser, msg *Message.Message){},
		responeMsg:          map[int64]chan *Message.Message{},
		cleanMemoryFragment: 1000,
		AuthenticationText:  []byte(defaultAuthenticationText),
	}
	fmt.Println(remoteAdderss)
	cli.ip, cli.prot, _ = utils.ParseAdderss(remoteAdderss)

	return &cli
}

func (c *Client) Connect(AuthenticationText ...string) error {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var err error

	if c.conn, err = net.DialTCP("tcp", &net.TCPAddr{net.IP{0, 0, 0, 0}, r.Intn(30000) + 2000, ""}, &net.TCPAddr{c.ip, c.prot, ""}); err != nil {
		return err
	}

	go c.read()

	if AuthenticationText != nil {
		c.AuthenticationText = []byte(AuthenticationText[0])
	}

	msg := Message.NewMsg(c.srcKey)
	msg.SetRequestString("Authentication", "Authentication", string(c.AuthenticationText))

	res, err := c.Request(msg)
	if res != nil && res.Body.StateCode == 0 {
		return errors.New(string(res.Body.Data))
	}

	return err
}

func (c *Client) read() {
	var msg *Message.Message
	var err error

	for {

		msg, err = c.UnMarshalMsg(c.conn)

		if err != nil {
			log.Println("client 解包失败连接断开", err)
			break
		}

		switch msg.Header.MType {
		// 响应请求
		case 0:
			hand, ok := c.HandlerFunc[msg.Header.DistApi]
			if !ok {
				msg.SetResponse(303, code[303])
				c.Response(msg)
				continue
			}
			go hand(c, msg)
		// 请求的响应
		case 1:
			res, ok := c.responeMsg[msg.Header.SrcApi]
			if !ok {
				c.pushMsg(msg)
				continue
			}
			res <- msg
		}

	}

	c.conn.Close()
}

func (c *Client) Request(message *Message.Message, timeOut ...time.Duration) (*Message.Message, error) {

	if timeOut == nil {
		timeOut = []time.Duration{c.TimeOut}
	}
	var cReq = make(chan *Message.Message)

	c.responeMsg[message.Header.SrcApi] = cReq

	mBuf, err := c.MarshalMsg(message)

	if err != nil {
		return nil, err
	}

	if _, err = c.conn.Write(mBuf); err != nil {
		close(cReq)
		return nil, errors.New("发送消息失败 " + err.Error())
	}

	var res *Message.Message

	select {

	case res = <-c.responeMsg[message.Header.SrcApi]:
		return res, nil
	case <-time.After(time.Second * timeOut[0]):
		return nil, errors.New(" request time out 请求超时")
	}

}

func (c *Client) Response(message *Message.Message) {

	buf, err := c.MarshalMsg(message)
	if err != nil {
		log.Println("client MarshalMsg err", err)
		return
	}
	_, err = c.conn.Write(buf)

	if err != nil {
		log.Println("client Write err", err)
		return
	}
}

func (c *Client) UnMarshalMsg(r io.Reader) (*Message.Message, error) {
	return utils.UnMarshalMsg(r)
}

func (c *Client) MarshalMsg(msg *Message.Message) ([]byte, error) {
	return utils.Marshal(msg)
}

func (c *Client) AddHandlerFunc(api string, handle func(res Responser, msg *Message.Message)) {
	if api == "" {
		log.Println("AddHandlerFunc err：api不能为空！")
		return
	}

	_, ok := c.HandlerFunc[api]
	if ok {
		log.Println("AddHandlerFunc err：api已存在！")
		return
	}
	c.HandlerFunc[api] = handle
}
func (c *Client) pushMsg(msg *Message.Message) {

	v := c.PushHanderFun

	if v == nil {
		log.Println("push message:", msg.String())
		return
	}

	v(msg)

}

//func (c *Client) request(distKey,distApi string ,data []byte) {
//	msg := Message.NewRequestMsg(c.srcKey,distKey,distApi,data)
//
//}
