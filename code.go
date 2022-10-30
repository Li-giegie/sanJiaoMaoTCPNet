package sanJiaoMaoTCPNet

var code = make(map[int][]byte)

func init()  {
	code[0] = []byte("认证失败")

	code[200] = []byte("server响应")
	code[201] = []byte("server响应：未找到响应API接口")

	code[300] = []byte("转发到目的主机处理")
	code[301] = []byte("转发失败 ：转发到目的主机失败 未找到目的主机")
	code[302] = []byte("转发失败 ：向目标主机发送消息失败")
	code[303] = []byte("响应成功 ：目标主机没有指定的API接口处理函数")


}
