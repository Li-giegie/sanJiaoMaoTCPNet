package Message

import "google.golang.org/protobuf/proto"

type MessageI interface {
	GetHeaderLen() uint32
	GetHeaderBytes() []byte
	GetHeader() Header
	GetDataLen() uint32
	GetDataBytes() []byte
}

type Message struct {
	Header
	Data []byte
}

func (m *Message) GetHeaderLen() uint32 {
	buf, err := proto.Marshal(&m.Header)
	if err != nil {
		return 0
	}
	return uint32(len(buf))
}

func (m *Message) GetHeaderBytes() []byte {
	res, _ := proto.Marshal(&m.Header)
	return res
}

func (m *Message) GetDataLen() uint32 {
	return uint32(len(m.Data))
}

func (m *Message) GetDataBytes() []byte {
	return m.Data
}

func (m Message) GetHeader() Header {
	return m.Header
}
