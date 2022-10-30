package Message

import (
	"encoding/json"
	"time"
)

// msgState:请求的的消息类型
func NewMsg(sk string) *Message {
	msg := Message{
		Header: &Header{SrcApi: time.Now().UnixNano(),SrcKey: sk},
		Body: &Body{StateCode: 0,Data: []byte{}},
	}
	return &msg
}

func (x *Message) SetResponseString(stateCode int,text string)  {
	x.SetResponse(stateCode,[]byte(text))
}

func (x *Message) SetResponseJSON(stateCode int,obj interface{})  {
	buf,_ :=json.Marshal(obj)
	 x.SetResponse(stateCode,buf)
}

func (x *Message) SetResponse(stateCode int,data []byte)  {
	x.Body.StateCode = int32(stateCode)
	x.Body.Data = data
	x.Header.MType = 1
}

func (x *Message) SetRequestString(dk,di,data string)  {
	x.SetRequest(dk,di,[]byte(data))
}

func (x *Message) SetRequestJson(dk,di string,data interface{})  {
	buf,_ := json.Marshal(data)
	x.SetRequest(dk,di,buf)
}


func (x *Message) SetRequest(dk,di string,data []byte)  {
	x.Header.MType,x.Body.StateCode = 0,200
	x.Header.DistKey,x.Header.DistApi = dk,di
	x.Body.Data = data
}
