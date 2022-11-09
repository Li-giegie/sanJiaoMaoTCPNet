package Message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
)

//打包方式二
func Pack(msg MessageI) ([]byte, error) {
	var msgBuf = new(bytes.Buffer)
	var err error
	hb := msg.GetHeaderBytes()
	// 写入消息头的长度
	err = binary.Write(msgBuf, binary.LittleEndian, uint32(len(hb)))
	if err != nil {
		return nil, errors.New("pack error -1：" + err.Error())
	}
	// 写入消息头
	_, err = msgBuf.Write(hb)
	if err != nil {
		return nil, errors.New("pack error -2：" + err.Error())
	}

	_, err = msgBuf.Write(msg.GetDataBytes())
	if err != nil {
		return nil, errors.New("pack error -3：" + err.Error())
	}
	return msgBuf.Bytes(), nil
}

func UnPack(r io.Reader) (*Message, error) {
	var headerLenb = make([]byte, 4)
	var err error

	var msg Message

	_, err = io.ReadFull(r, headerLenb)

	if err != nil {
		return nil, errors.New("UnPack err -1:" + err.Error())
	}
	// header长度
	headerLen_uint32 := binary.LittleEndian.Uint32(headerLenb)
	// header protoc 字节码
	var headerBytes = make([]byte, headerLen_uint32)
	_, err = io.ReadFull(r, headerBytes)
	if err != nil {
		return nil, errors.New("UnPack err -2" + err.Error())
	}

	err = proto.Unmarshal(headerBytes, &msg.Header)
	if err != nil {
		return nil, errors.New("UnPack err -3" + err.Error())
	}

	msg.Data = make([]byte, msg.Header.DataLen)
	_, err = io.ReadFull(r, msg.Data)
	if err != nil {
		return nil, errors.New("UnPack err -4" + err.Error())
	}

	return &msg, nil
}
