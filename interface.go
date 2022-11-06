package sanJiaoMaoTCPNet

import (
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"time"
)

type ServerI interface {
	// 设置客户端认证回调 可选 返回值决定认证是否成功 err 为空是认证通过 否则与客户端断开连接
	AddAuthenticationHandle(fun func(key, ip string, text []byte) error)
	// 运行
	Run() error
	// 添加接口 处理回调
	AddHandlerFunc(api string, handle func(msg *Message.Message))
	// 设置服务端最大容忍连接错误数量 值 >0 时有意义
	SetMaxErrorConnectNum(num int)
	// 停止侦听 并释放所有建立的连接
	Shutdown()
}

type ClientI interface {
	// 设置认证文本
	SetAuthenticationText(AuthenticationText []byte)
	// 设置本地地址端口号可选项 默认0.0.0.0:10000 ~ 65535
	SetLocalAdderss(address string) error
	// 与服务器建立连接 参数为可选 认证文本 设置过忽略此方法
	Connect(AuthenticationText ...string) error

	// 客户端作为请求端时次函数 为请求放 参数 消息内容 超时时间
	Request(message *Message.Message, timeOut ...time.Duration) (*Message.Message, error)
	// 客户端作为响应端是应用此函数 作为回复函数
	Response(message *Message.Message)
	// 以接口形式添加 处理方法 参数: res 作为响应端  msg 为消息
	AddHandlerFunc(api string, handle func(res Responser, msg *Message.Message))
	// 服务端可能会推送消息 此方法为推送消息处理回调
	AddPushHandlerFunc(fun func(msg *Message.Message))

	Close()
}
