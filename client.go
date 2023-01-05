package sanJiaoMaoTCPNet

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type ClientI interface {
	SetDefaultSendKey(key string)
	AddHandlerFunc(api string, handler func(msg *Message.Message, res Message.ReplyMessageI))
	//SetPushMessageHandler(pushHandlerFunc func(msg *Message.MessageType2,res utils.ReplyMessageI))
	Connect(authentication ...string) error
	SendMessage(distKey, distApi string, stateCode int, data []byte, timeOut ...time.Duration) (*Message.Message, error)
	SendMessageString(distKey, distApi string, stateCode int, data string, timeOut ...time.Duration) (*Message.Message, error)
	SendMessageJSON(distKey, distApi string, stateCode int, data interface{}, timeOut ...time.Duration) (*Message.Message, error)
	Close()
	Run()
}
type handlerFunc map[string]func(msg *Message.Message, res Message.ReplyMessageI)
type Client struct {
	key                string
	laddr              *net.TCPAddr
	raddr              *net.TCPAddr
	conn               *net.TCPConn
	TimeOut            time.Duration
	replyChan          map[int64]chan *Message.Message
	handlerFunc        map[string]func(msg *Message.Message, res Message.ReplyMessageI)
	defaultSendKey     string
	AutoId             int64
	pushMessageHandler func(msg *Message.Message, res Message.ReplyMessageI)
	isClose            bool
	mutex              sync.RWMutex


}

func NewClient(raddr string, key string) ClientI {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var client = Client{key: key,
		TimeOut:     time.Second * 20,
		replyChan:   map[int64]chan *Message.Message{},
		handlerFunc: map[string]func(msg *Message.Message, res Message.ReplyMessageI){},
	}
	client.laddr, _ = net.ResolveTCPAddr("tcp", "0.0.0.0:"+fmt.Sprintf("%v", r.Intn(60000)+2000))

	client.raddr, _ = net.ResolveTCPAddr("tcp", raddr)

	return &client
}

func (c *Client) Connect(authentication ...string) error {
	var err error
	c.conn, err = net.DialTCP("tcp", c.laddr, c.raddr)
	if err != nil {
		return err
	}

	if authentication == nil {
		authentication = []string{""}
	}

	go c.read()

	reply, err := c.SendMessageString("", "", 0, authentication[0])

	if err != nil {
		c.Close()
		return err
	}
	if reply.StateCode == 0 {
		c.Close()
		return errors.New(string(reply.Data))
	}

	return nil
}

func (c *Client) read() {

	for {
		msg, err := Message.UnPack(c.conn)
		if err != nil {
			log.Println("client read msg err：", err)
			return
		}
		// 请求
		//fmt.Println(msg.SrcApi,msg.String(),string(msg.Data))
		switch msg.MType {
		case 0:
			reply := Message.NewReplyMsg(msg, c.conn)
			handler, ok := c.handlerFunc[msg.DistApi]
			if !ok {
				err = reply.String(301, "在目标客户端找不到处理函数！")
				fmt.Println(err)
				continue
			}
			go handler(msg, reply)
		case 1:
			c.mutex.RLock()
			reply, ok := c.replyChan[msg.SrcApi]

			if !ok  {
				fmt.Println("push message:", msg.String(), string(msg.Data))
				continue
			}
			reply <- msg
			c.mutex.RUnlock()

		default:
			log.Println("未知消息类型")
		}

	}
}

func (c *Client) SendMessage(distKey, distApi string, stateCode int, data []byte, timeOut ...time.Duration) (*Message.Message, error) {
	if distKey == "" {
		distKey = c.defaultSendKey
	}
	if timeOut == nil {
		timeOut = []time.Duration{c.TimeOut}
	}

	c.mutex.Lock()
	c.AutoId++

	var msg = &Message.Message{
		Header: Message.Header{
			SrcKey:    c.key,
			SrcApi:    c.AutoId,
			DistKey:   distKey,
			DistApi:   distApi,
			StateCode: int32(stateCode),
			DataLen:   uint32(len(data)),
		},
		Data: data,
	}

	var reply = make(chan *Message.Message)

	c.replyChan[msg.SrcApi] = reply
	c.mutex.Unlock()

	buf, err := Message.Pack(msg)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write(buf)
	if err != nil {
		return nil, err
	}

	select {
	case res := <-reply:
		close(reply)
		return res, nil
	case <-time.After(timeOut[0]):
		close(reply)
		return nil, errors.New("timeOut SrcApi: " + strconv.Itoa(int(msg.SrcApi)))
	}

}

func (c *Client) SendMessageString(distKey, distApi string, stateCode int, data string, timeOut ...time.Duration) (*Message.Message, error) {
	return c.SendMessage(distKey, distApi, stateCode, []byte(data), timeOut...)
}

func (c *Client) SendMessageJSON(distKey, distApi string, stateCode int, data interface{}, timeOut ...time.Duration) (*Message.Message, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.SendMessage(distKey, distApi, stateCode, buf, timeOut...)
}

func (c *Client) AddHandlerFunc(api string, handler func(msg *Message.Message, res Message.ReplyMessageI)) {
	_, ok := c.handlerFunc[api]
	if !ok {
		c.handlerFunc[api] = handler
		return
	}
	log.Println("warning api exits")
}

func (c *Client) SetDefaultSendKey(key string) { c.defaultSendKey = key }

func (c *Client) Run() {
	for !c.isClose {
	}
}

func (c *Client) Close() {

	if c.conn != nil {
		c.conn.Close()
		c.conn.CloseRead()
		c.conn.CloseWrite()
	}
	c.isClose = true
}
