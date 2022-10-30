package sanJiaoMaoTCPNet

import (
	"errors"
	"io"
	"log"
	"net"
	"sanJiaoMaoTCPNet/Message"
	"sanJiaoMaoTCPNet/utils"
	"time"
)

const (
	MSGSTATE_CLOSE byte= 0
	MSGSTATe_KEEP byte= 1
)

type Server struct {

	Key string
	ip net.IP
	port int

	aclList map[string] aclList
	conns map[string] connections

	MaxErrorConnectNum int
	AuthenticationHandle func(key,ip string,text []byte) error

	HandlerFunc map[string] func(msg *Message.Message)

}

type connections struct {
	conn *net.TCPConn
}

type aclList struct {
	id int64
	errorConnectNum int

}

func NewServer(adderss,Key string) *Server  {

	ip,port,_ := utils.ParseAdderss(adderss)
	return &Server{
		Key:                  Key,
		ip:                   ip,
		port:                 port,
		MaxErrorConnectNum:   3,
		conns: map[string]connections{},
		aclList: map[string]aclList{},
		HandlerFunc: map[string]func(msg *Message.Message){},
		AuthenticationHandle: nil,
	}
}

func (s *Server) Run() error {

	listen,err := net.ListenTCP("tcp",&net.TCPAddr{s.ip,s.port,""})

	if err != nil { return err }

	defer listen.Close()

	for  {
		conn,err :=listen.AcceptTCP()

		if err !=nil { continue }

		go s.process(conn)
	}

}

func (s *Server) authentication(msg *Message.Message,ip string) error {

	// 1.检查是否拒绝登录列表
	v,ok := s.aclList[ip]
	if !ok {
		v = aclList{time.Now().UnixNano(), 0}
		s.aclList[ip] = v
	}
	if v.errorConnectNum >= s.MaxErrorConnectNum { return errors.New("服务器拒绝建立连接 ")}

	// 2.检查是否上线
	_,cok := s.conns[msg.Header.SrcKey]
	if cok  {
		return errors.New("非法认证 认证失败 此用户以上线")
	}

	// 3.无认证过程直接通过 有认证过程执行自定义认证过程
	func1 := s.AuthenticationHandle
	if func1 == nil {
		return nil
	}

	return func1(msg.Header.SrcKey,ip,msg.Body.Data)

}

func (s *Server) process(conn *net.TCPConn)  {
	var msg *Message.Message
	var err error
	var ip,key string

	if msg,err = s.UnMarshalMsg(conn);err != nil {
		log.Println("解析认证消息错误 程序结束！")
		conn.Close()
		delete(s.conns,key)
		return
	}

	ip,key = conn.RemoteAddr().String(),msg.Header.SrcKey

	v:= s.aclList[ip];v.id=time.Now().UnixNano()

	// 认证失败
	if err = s.authentication(msg,ip); err != nil {
		v.errorConnectNum ++;s.aclList[ip] = v
		msg.SetResponseString(0,err.Error())
		s.write(conn,msg)
		return
	}

	msg.SetResponseString(1,"success")
	err = s.write(conn,msg)
	if err != nil { return }

	// 上线
	v.errorConnectNum = 0
	s.conns[msg.Header.SrcKey],s.aclList[ip] = connections{ conn},v

	var conns connections
	var ok bool
	var tempKey string
	for  {
		msg,err = s.UnMarshalMsg(conn)
		if err != nil { break }

		// 服务器处理函数
		if msg.Header.DistKey == s.Key {

			sf,ok := s.HandlerFunc[msg.Header.DistApi]
			if !ok {
				msg.SetResponse(201,code[201])
				s.write(conn,msg)
				continue
			}
			sf(msg)
			s.write(conn,msg)
			continue
		}

		if msg.Header.MType == 0 {
			tempKey = msg.Header.DistKey
		}else {
			tempKey = msg.Header.SrcKey
		}

		// 转发到其他客户端
		conns,ok = s.conns[tempKey]
		if !ok {
			msg.SetResponse(301,code[301])
			s.write(conn,msg)
			continue
		}
			//转发
		err = s.write(conns.conn,msg)
		if err != nil {
			msg.SetResponse(302,code[302])
			s.write(conn,msg)
		}


	}

	conn.Close()
	log.Println("off line :key---",key,"ip---",ip)
	delete(s.conns,key)
}

func (s *Server) write(conn *net.TCPConn,msg *Message.Message) error {

	msgBuf,err := s.MarshalMsg(msg)

	if err != nil {
		return errors.New("服务端序列化消息失败 :"+err.Error())
	}

	_,err = conn.Write(msgBuf)

	if err != nil { log.Println("server write err:",err ) }

	if msg.Body.StateCode == 0 {
		conn.Close()
	}
	return err
}

func (c *Server) UnMarshalMsg(r io.Reader)(*Message.Message,error)  {
	return utils.UnMarshalMsg(r)
}

func (c *Server) MarshalMsg(msg *Message.Message) ([]byte,error) {
	return utils.Marshal(msg)
}