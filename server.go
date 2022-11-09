package sanJiaoMaoTCPNet

import (
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
	"net"
)

type ServerI interface {
	AddHandleFunc(api string, handle func(msg *Message.Message, reply Message.ReplyMessageI))
	SetAuthentication(authentication func(ip string, key string, data []byte) (bool, string))
	Run()
	Shutdown()
}

type Server struct {
	key            string
	addr           *net.TCPAddr
	listen         *net.TCPListener
	err            error
	conns          map[string]*Connecter
	Authentication func(ip string, key string, data []byte) (bool, string)
	handleFunc     map[string]func(msg *Message.Message, reply Message.ReplyMessageI)
	isShutdown     bool
}

// 管理的连接对象
type Connecter struct {
	ip    string
	conn  *net.TCPConn
	key   string
	state int
}

func NewServer(addr string, key ...string) (ServerI, error) {
	var s = &Server{
		isShutdown: false,
		conns:      map[string]*Connecter{},
		handleFunc: map[string]func(msg *Message.Message, reply Message.ReplyMessageI){},
	}

	s.addr, s.err = net.ResolveTCPAddr("tcp", addr)
	if s.err != nil {
		return nil, s.err
	}

	if key == nil {
		key = []string{"sanjiaomao"}
	}
	s.key = key[0]
	s.listen, s.err = net.ListenTCP("tcp", s.addr)
	if s.err != nil {
		return nil, s.err
	}
	return s, nil
}

func (s *Server) Run() {
	for {

		if s.isShutdown {
			if s.listen != nil {
				s.listen.Close()
			}
			for _, connecter := range s.conns {
				if connecter.conn != nil {
					connecter.conn.Close()
				}
				delete(s.conns, connecter.key)
			}
			break
		}

		conn, err := s.listen.AcceptTCP()
		if err != nil {
			log.Println(err)
			break
		}
		go s.Read(conn)
	}
	log.Println("shutdown ----------")
}

func (s *Server) Read(conn *net.TCPConn) {
	var info = "Authentication success"

	msg, err := Message.UnPack(conn)
	var key = msg.SrcKey
	if err != nil {
		log.Println("解包失败：", err)
		conn.Close()
		return
	}
	reply := Message.NewReplyMsg(msg, conn)

	//根据key查看是否已经认证过
	_, ok := s.conns[msg.SrcKey]
	if ok {
		if err = reply.String(0, "拒绝认证：线上已有用户"); err != nil {
			log.Println("Authentication reply string :", err)
		}
		conn.Close()
		return
	}
	if auth := s.Authentication; auth != nil {

		ok, info = auth(conn.RemoteAddr().String(), msg.SrcKey, msg.Data)
		if !ok {
			if err = reply.String(0, info); err != nil {
				log.Println("Authentication reply string :", err)
			}
			conn.Close()
			return
		}
	}
	// 认证通过 回复一个消息 失败断开连接
	if err = reply.String(1, info); err != nil {
		conn.Close()
		return
	}

	s.conns[msg.SrcKey] = &Connecter{
		ip:    conn.RemoteAddr().String(),
		conn:  conn,
		key:   msg.DistKey,
		state: 1,
	}

	fmt.Println("认证一个:", s.conns, "key = ", msg.SrcKey)
	for {
		msg, err = Message.UnPack(conn)

		if err != nil {
			log.Println("解包失败：", err)
			break
		}
		// 服务端处理API
		if msg.DistKey == s.key {
			reply.SetMsg(msg)
			handler, ok := s.handleFunc[msg.DistApi]
			if !ok {
				reply.String(101, "服务端找不到处理函数！")
				continue
			}
			go handler(msg, reply)
		} else { //转发消息
			temKey := msg.DistKey
			if msg.MType == 1 {
				temKey = msg.SrcKey
			}
			conner, ok := s.conns[temKey]
			if !ok {
				reply.SetMsg(msg)
				err = reply.String(201, "转发消息失败 原因是目的客户端不存在！")
				if err != nil {
					log.Println(201, "转发消息失败 原因是目的客户端不存在！", err)
				}
				continue
			}
			go s.forward(conn, conner.conn, msg)
		}

	}

	s.CLoseConn(key)
}

func (s *Server) SetAuthentication(authentication func(ip string, key string, data []byte) (bool, string)) {
	s.Authentication = authentication
}

// conn is src conn forwardconn is dist conn
func (s *Server) forward(conn *net.TCPConn, forwardConn *net.TCPConn, msg *Message.Message) {
	buf, err := Message.Pack(msg)
	if err != nil {
		if err = Message.NewReplyMsg(msg, conn).String(202, "forward Pack err："+err.Error()); err != nil {
			log.Println("NewReplyMsg forward err", err)
			s.CLoseConn(msg.SrcKey)
		}
		return
	}
	_, err = forwardConn.Write(buf)
	if err != nil {
		if err = Message.NewReplyMsg(msg, conn).String(203, "forward err："+err.Error()); err != nil {
			s.CLoseConn(msg.SrcKey)
		}
		log.Println("forward err:", err)
	}
}

func (s *Server) CLoseConn(key string) {
	if s.conns[key].conn != nil {
		s.conns[key].conn.Close()
	}

	delete(s.conns, key)

}

func (s *Server) AddHandleFunc(api string, handle func(msg *Message.Message, reply Message.ReplyMessageI)) {
	_, ok := s.handleFunc[api]
	if ok {
		log.Println("AddHandleFunc err:api exits")
		return
	}
	s.handleFunc[api] = handle
}

func (s Server) Shutdown() {
	s.isShutdown = true
}
