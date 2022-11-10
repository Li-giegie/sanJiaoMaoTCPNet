​        [简体中文](#简体中文)

​        [English](#English)

## English

# A Simple Concurrent TCP Framework Based on Golang protobuf Message Header

## Support the client to request the server, just like HTTP, based on the interface

## You can also forward the public network request to the intranet client and write your own processing program

## It supports forwarding requests to other clients connected to the server

### For example, if you have a public network server, you can deploy the server in your public network server. If you start two clients in the intranet, you can communicate.

### 1、Basic functions

    (1).Support the server to forward requests
    
    (2).Support client side parallel requests
    
    (3).Support the server interface to process requests

    (4).Support authentication before connection

### 2、Core functions

    (1).client -> server The client sends a request to the server for processing
    
    (2).client1 -> server -> client2  Client 1 sends it to the server and the server forwards it to client 2

### 3、test data

    (1).testing environment  
    Windows 12 
    Notebook computer 
    Video card GTX1650
    CPU	AMD Ryzen 5 5600H with Radeon Graphics            3.30 GHz
    Memory RAM	16.0 GB (13.9 GB available)
    System Type 64 bit operating system, X64 based processor

---

    (2).Client requests server Local environment 1024Byte 10000 requests time 1.714s

![img_1.png](img_1.png)

---

    3.Client requests server forwarding Local environment 1024Byte 10000 requests time 1.1714s

![img_2.png](img_2.png)

### 4、Server Example

    A simple server only needs three
    (1). Create server object
    Two parameters must be passed: the server's listening address and the Key used by the server to be found by the client
    srv, err := NewServer("127.0.0.1:9999", "server")

---

    (2).Add a handler
    srv.AddHandleFunc("test", func(msg *Message.Message, reply Message.ReplyMessageI) {
		//fmt.Println("test----", msg.String(), string(msg.Data))
		reply.String(200, "server 收到")
    })

___

    (3).Start it
    srv.Run()

#### Tips

The full version is displayed on the server_ test. In go TestServer, there is an additional SetAuthentication method
compared with the above content. This method is used to authenticate the client before establishing a communication
connection with the server. A callback method needs to be passed,

The return value of the callback method (is_pass bool, info string), the return value of the boolean value, determines
whether to establish communication with the client. The string is used as the reply information. If this method is not
set, communication will be established by default.
___

### 5.Client Example

    leave out the nonessential words
    //Create a new client parameter: 1. Remote connection address; 2. Customizing a unique client ID
    c = NewClient("127.0.0.1:9999", "client1")
    
    // Connect the optional parameter authentication text. The authentication text will be sent to the server for verification. If the verification fails, the connection establishment fails. Please ensure that the server is authorized

    err := c.Connect()
    // Try sending a message
    Parameter 1: key identifier of the server or client of the remote connection, 2. interface of the connection, 3. status code, 3. message content, 4. optional timeout
    res, err := c.SendMessage("client2", "test", 200, testBuf, time.Second*2)

## 简体中文

# 一个简基于Golang protobuf消息头 的简单 并发TCP 框架

## 支持客户端请求服务端 就像HTTP那样 基于接口的形式

## 也可以把公网的请求转发到内网客户端上面自己写处理程序

## 支持把请求转发到与服务端建立连接的其他客户端上

### 例如你拥有一台公网服务器可以把服务端部署在你的公网服务器中，内网在启动两个客户端就可以通信了。

### 1、基本的功能

    (1).支持服务端转发请求
    
    (2).支持客户端并行请求
    
    (3).支持服务端接口处理请求 
    
    (4).支持建立连接前的认证

### 2、核心功能

    (1).client -> server 客户端发起请求服务端处理
    
    (2).client1 -> server -> client2  客户端1 发给服务端 服务端转发到 客户端2

### 3、测试数据

    (1).测试环境  
    Windows 12 
    笔记本电脑 
    显卡 1650 内存 13.9
    处理器	AMD Ryzen 5 5600H with Radeon Graphics            3.30 GHz
    机带 RAM	16.0 GB (13.9 GB 可用)
    系统类型	64 位操作系统, 基于 x64 的处理器
    笔和触控	没有可用于此显示器的笔或触控输入

---

    (2).客户端请求服务端 本地环境 1024Byte 10000次请求 耗时1.714s

![img_1.png](img_1.png)

---

    3.客户端请求服务器转发 本地环境 1024Byte 10000次请求 耗时1.1714s

![img_2.png](img_2.png)

### 4、服务端例子

    一个简单的服务端 仅仅需要三部
    (1). 创建服务端对象
    必须传递两个参数 服务端侦听地址 和服务端用于被客户端发现的Key key为自定义内容
    srv, err := NewServer("127.0.0.1:9999", "server")

---

    (2).添加一个处理函数
    srv.AddHandleFunc("test", func(msg *Message.Message, reply Message.ReplyMessageI) {
		//fmt.Println("test----", msg.String(), string(msg.Data))
		reply.String(200, "server 收到")
    })

___

    (3).启动它
    srv.Run()

#### 提示

完整版的内容展现在server_test.go TestServer 中，与上述内容相比多出了一个 SetAuthentication的方法 ，此方法为客户端与服务端建立通信连接前的认证，需传递一个回调方法，
回调方法的返回值（is_pass bool,info string）布尔值的返回值将决定是否与客户端建立通信， 字符串作为回复信息，此方法不设置将默认建立通信。
___

### 5.客户端例子

    闲言少叙
    //新建一个客户端 参数 1 远程连接的地址，2.一个唯一的 客户端标识 自定义 
    c = NewClient("127.0.0.1:9999", "client1")
    
    // 连接 可选参数 认证文本 ，认证文本将发送至服务端进行校验，校验失败建立连接失败，请确保服务端授权
    err := c.Connect()
    // 发一条消息 试试
    参数1 远程连接的服务端或客户端 key标识，2.连接的 接口，3.状态码，3.消息内容，4.可选的超时时间
    res, err := c.SendMessage("client2", "test", 200, testBuf, time.Second*2)
