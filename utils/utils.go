package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func ParseAdderss(adderss string) (net.IP, int, error) {
	readderss := strings.Split(adderss, ":")

	if len(readderss) != 2 {
		return net.IP{}, 0, errors.New("地址格式错误")
	}
	rePort, err := strconv.Atoi(readderss[1])

	if err != nil {
		return net.IP{}, 0, errors.New("地址格式错误")
	}

	var ip = net.ParseIP(readderss[0])

	if ip == nil {
		return net.IP{}, 0, errors.New("地址格式错误")
	}

	return ip, rePort, nil
}

func ReadMessage(r io.Reader) ([]byte, error) {
	var mlen = make([]byte, 4)
	var mlen_uint32 uint32
	var err error

	_, err = io.ReadFull(r, mlen)
	if err != nil {
		return nil, err
	}

	mlen_uint32 = binary.LittleEndian.Uint32(mlen)

	var msgBuf = make([]byte, mlen_uint32)

	_, err = io.ReadFull(r, msgBuf)
	if err != nil {
		return nil, err
	}

	return msgBuf, err
}
func UnMarshalMsg(r io.Reader) (*Message.Message, error) {

	var msg Message.Message

	msgBuf, err := ReadMessage(r)

	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(msgBuf, &msg)

	return &msg, err
}

func Marshal(msg *Message.Message) ([]byte, error) {
	var mbuf = new(bytes.Buffer)
	var err error

	var msgb []byte

	msgb, err = proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	err = binary.Write(mbuf, binary.LittleEndian, uint32(len(msgb)))
	if err != nil {
		return nil, err
	}

	_, err = mbuf.Write(msgb)
	if err != nil {
		return nil, err
	}

	return mbuf.Bytes(), nil

}

func Input(cmdtipchar ...string) []string {
	if cmdtipchar == nil {
		cmdtipchar = []string{">>"}
	}
	inputReader := bufio.NewReader(os.Stdin) //创建一个读取器，并将其与标准输入绑定。
	var inputs string
	var err error
	var cmm []string
	for {
		fmt.Printf(cmdtipchar[0])
		inputs, err = inputReader.ReadString('\n') //读取器对象提供一个方法 ReadString(delim byte) ，该方法从输入中读取内容，直到碰到 delim 指定的字符，然后将读取到的内容连同 delim 字符一起放到缓冲区。

		if err != nil {
			fmt.Println("运行异常：", err)
			continue
		}

		if inputs == "" {
			continue
		}
		inputs = strings.Replace(inputs, "\r", "", -1)
		inputs = strings.Replace(inputs, "\n", "", -1)

		for _, s := range strings.Split(inputs, " ") {
			if s == "" || s == " " {
				continue
			}
			cmm = append(cmm, s)
		}
		if cmm != nil && len(cmm) > 0 {
			break
		}
	}

	return cmm
}
