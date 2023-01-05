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

func PackV2(msg *Message) ([]byte,error) {
	var msgBuf = new(bytes.Buffer)
	var err error

	var packBuf = new(bytes.Buffer)

	packBuf.WriteString(msg.SrcKey+"\r\n" + msg.DistKey+"\r\n" + msg.DistApi+"\r\r\n")


	err = binary.Write(packBuf, binary.LittleEndian, msg.SrcApi)
	if err != nil {return nil, err}
	err = binary.Write(packBuf, binary.LittleEndian, msg.MType)
	if err != nil {return nil, err}
	err = binary.Write(packBuf, binary.LittleEndian, msg.StateCode)
	if err != nil {return nil, err}
	err = binary.Write(packBuf, binary.LittleEndian, msg.DataLen)
	if err != nil {return nil, err}

	// 写入消息头的长度
	err = binary.Write(msgBuf, binary.LittleEndian, uint32(packBuf.Len() +len(msg.Data)))
	if err != nil {return nil, err}
	if err != nil {
		return nil, errors.New("pack error -1：" + err.Error())
	}
	// 写入消息头
	_, err = msgBuf.Write(packBuf.Bytes())
	if err != nil {
		return nil, errors.New("pack error -2：" + err.Error())
	}
	_,err = msgBuf.Write(msg.Data)
	if err != nil {
		return nil, errors.New("pack error -2：" + err.Error())
	}
	return msgBuf.Bytes(), nil
}

func UnPackV2(r io.Reader) (*Message,error)  {
	var dataLen = make([]byte, 4)
	var err error

	_, err = io.ReadFull(r, dataLen)

	if err != nil {
		return nil, errors.New("UnPack err -1:" + err.Error())
	}

	// header长度
	dataLen_uint32 := binary.LittleEndian.Uint32(dataLen)

	// header protoc 字节码
	var dataBytes = make([]byte, dataLen_uint32)

	_, err = io.ReadFull(r, dataBytes)
	if err != nil {
		return nil, errors.New("UnPack err -2" + err.Error())
	}

	var msg Message

	//smsd
	i := bytes.Index(dataBytes,[]byte("\r\r\n"))
	q:= bytes.Split(dataBytes[:i],[]byte("\r\n"))
	msg.SrcKey = string(q[0])
	msg.DistKey = string(q[1])
	msg.DistApi = string(q[2])
	i+=3
	msg.SrcApi = int64(binary.LittleEndian.Uint64(dataBytes[i:i+8]))
	i+=8
	msg.MType = int32(binary.LittleEndian.Uint32(dataBytes[i:i+4]))
	i+=4
	msg.StateCode = int32(binary.LittleEndian.Uint32(dataBytes[i:i+4]))
	i+=4
	msg.StateCode = int32(binary.LittleEndian.Uint32(dataBytes[i:i+4]))
	i+=4

	msg.Data = dataBytes[i:]

	return &msg, nil
}