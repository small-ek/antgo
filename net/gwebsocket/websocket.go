package gwebsocket

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/small-ek/ginp/crypto/hash"
	"sync"
)

//Connection ...
type Connection struct {
	Address       string          // 客户端地址
	AppId         string          // 登录的平台Id app/web/ios
	UserId        string          // 用户标识，用户登录以后才有
	LoginTime     uint64          // 登录时间 登录以后才有
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	TimeOut       uint64          // 超时时间
	Socket        *websocket.Conn // 用户连接
	ReadChan      chan []byte     // 读取消息通道
	WriteChan     chan []byte     // 写入消息发送客户端
	CloseChan     chan byte       // 关闭通道
	mutex         sync.Mutex      // 对关闭上锁
	isClosed      bool            // 防止closeChan被关闭多次
}

//New ...
func New(socket *websocket.Conn, address string, firstTime uint64) *Connection {
	var get = &Connection{
		Address:       address,
		Socket:        socket,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		TimeOut:       360,
		ReadChan:      make(chan []byte, 1000),
		WriteChan:     make(chan []byte, 1000),
		CloseChan:     make(chan byte, 1),
	}
	// 启动读协程
	go get.readLoop()
	// 启动写协程
	go get.writeLoop()
	return get
}

//ReadMessage Read client messages<读取客户端消息>
func (get *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-get.ReadChan:
	case <-get.CloseChan:
		err = errors.New("connection is closeed")
	}
	return
}

//WriteMessage Write message sending client(写入消息发送客户端)
func (get *Connection) WriteMessage(data []byte) (err error) {
	select {
	case get.WriteChan <- data:
	case <-get.CloseChan:
		err = errors.New("connection is closeed")
	}
	return
}

//Close Close the connection<关闭连接>
func (get *Connection) Close() {
	// 线程安全，可多次调用
	get.Socket.Close()
	// 利用标记，让closeChan只关闭一次
	get.mutex.Lock()
	if !get.isClosed {
		close(get.CloseChan)
		get.isClosed = true
	}
	defer get.mutex.Unlock()
}

//readLoop Read the channel message loop<读取通道消息循环>
func (get *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = get.Socket.ReadMessage(); err != nil {
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置
		select {
		case get.ReadChan <- data:
		case <-get.CloseChan: // closeChan 感知 conn断开
			goto ERR
		}

	}

ERR:
	get.Close()
}

//writeLoop Writes to the channel message loop<写入通道消息循环>
func (get *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-get.WriteChan:
		case <-get.CloseChan:
			goto ERR
		}
		if err = get.Socket.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	get.Close()
}

//SetHeartbeat 设置用户心跳
func (get *Connection) SetHeartbeat(currentTime uint64) {
	get.HeartbeatTime = currentTime
	return
}

//SetTimeOut 设置超时时间
func (get *Connection) SetTimeOut(TimeOut uint64) {
	get.TimeOut = TimeOut
	return
}

//GetHeartbeat 获取心跳是否超时
func (get *Connection) GetHeartbeat(currentTime uint64) bool {
	if get.HeartbeatTime+get.TimeOut <= currentTime {
		return true
	}
	return false
}

//SetLogin 设置用户登录
func (get *Connection) SetLogin(appId string, userId string, loginTime uint64) {
	get.AppId = appId
	get.UserId = userId
	get.LoginTime = loginTime
	get.SetHeartbeat(loginTime)
}

//GetLogin 获取是否登录
func (get *Connection) GetLogin() bool {
	if get.UserId != "" {
		return true
	}
	return false
}

//Login 用户登录
type Login struct {
	AppId  string
	UserId string
	Client *Connection
}

//GetUserKey 获取用户的Key
func (get *Login) GetUserKey() (key string) {
	key = hash.Sha1(get.AppId + get.UserId)
	return
}
