package Message

import (
	"encoding/json"
	"net"
)

type ReplyMessageI interface {
	String(stateCode int, data string) error
	Bytes(stateCode int, data []byte) error
	JSON(stateCode int, data interface{}) error
}

type ReplyMessage struct {
	msg  *Message
	conn *net.TCPConn
}

func NewReplyMsg(msg *Message, conn *net.TCPConn) *ReplyMessage {
	msg.MType = 1
	return &ReplyMessage{
		msg:  msg,
		conn: conn,
	}
}
func (r *ReplyMessage) String(stateCode int, data string) error {
	return r.Bytes(stateCode, []byte(data))
}

func (r *ReplyMessage) JSON(stateCode int, data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Bytes(stateCode, buf)
}

func (r *ReplyMessage) Bytes(stateCode int, data []byte) error {
	r.msg.Data = data
	r.msg.DataLen = uint32(len(r.msg.Data))
	r.msg.Header.StateCode = int32(stateCode)
	bur, err := Pack(r.msg)
	if err != nil {
		return err
	}
	_, err = r.conn.Write(bur)
	return err
}

func (r *ReplyMessage) SetMsg(msg *Message) {
	msg.MType = 1
	r.msg = msg
}
