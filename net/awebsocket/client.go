package awebsocket

import (
	"github.com/small-ek/antgo/crypto/hash"
	"sync"
)

//Client 客户端连接管理
type Client struct {
	Clients     map[*Connection]bool   // 全部的连接
	ClientsLock sync.RWMutex           // 读写锁
	Users       map[string]*Connection // 登录的用户
	UserLock    sync.RWMutex           // 读写锁
	Register    chan *Connection       // 连接通道处理
	Login       chan *Login            // 用户登录通道处理
	Close       chan *Connection       // 断开连接处理程序
	Broadcast   chan []byte            // 广播消息通道处理
}

//NewClient 默认初始化客户端
func NewClient() (clientManager *Client) {
	return &Client{
		Clients:   make(map[*Connection]bool),
		Users:     make(map[string]*Connection),
		Register:  make(chan *Connection, 1000),
		Login:     make(chan *Login, 1000),
		Close:     make(chan *Connection, 1000),
		Broadcast: make(chan []byte, 1000),
	}
}

//IsClient 是否存在
func (get *Client) IsClient(client *Connection) bool {
	get.ClientsLock.RLock()
	defer get.ClientsLock.RUnlock()
	return get.Clients[client]
}

//GetClient Get client 获取客户端连接
func (get *Client) GetClient() (clients map[*Connection]bool) {
	clients = make(map[*Connection]bool)
	get.GetClientsLoop(func(client *Connection, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

//GetClientsLoop Loop all connections<循环所有的客户端连接>
func (get *Client) GetClientsLoop(f func(client *Connection, value bool) (result bool)) {
	get.ClientsLock.RLock()
	for key, value := range get.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}
	defer get.ClientsLock.RUnlock()
	return
}

//GetClientsCount Gets the length of the connection<获取客户端的总长度>
func (get *Client) GetClientsCount() int {
	return len(get.Clients)
}

//AddClients Adding a connection 添加客户端
func (get *Client) AddClients(client *Connection) {
	get.ClientsLock.Lock()
	get.Clients[client] = true
	defer get.ClientsLock.Unlock()
}

//DeleteClients Delete client<删除客户端>
func (get *Client) DeleteClients(client *Connection) {
	get.ClientsLock.Lock()
	if _, ok := get.Clients[client]; ok {
		delete(get.Clients, client)
	}
	defer get.ClientsLock.Unlock()
}

//GetUserKey 获取用户key
func GetUserKey(appId string, userId string) string {
	return hash.Sha1(appId + userId)
}

//AddUsers 添加用户
func (get *Client) AddUsers(key string, connection *Connection) {
	get.UserLock.Lock()
	get.Users[key] = connection
	defer get.UserLock.Unlock()
}

//GetUserClient 获取用户的连接
func (get *Client) GetUserClient(appId string, userId string) (connection *Connection) {
	get.UserLock.RLock()
	userKey := GetUserKey(appId, userId)
	if value, ok := get.Users[userKey]; ok {
		connection = value
	}
	defer get.UserLock.RUnlock()
	return
}

//GetUsersCount Get the total number of users<获取用户总数>
func (get *Client) GetUsersCount() int {
	return len(get.Users)
}

//DeleteUsers Delete user<删除用户>
func (get *Client) DeleteUsers(connection *Connection) (result bool) {
	get.UserLock.Lock()
	key := GetUserKey(connection.AppId, connection.UserId)
	if value, ok := get.Users[key]; ok {
		// 判断是否为相同的用户
		if value.Address != connection.Address {
			return
		}
		delete(get.Users, key)
		result = true
	}
	defer get.UserLock.Unlock()
	return
}

//GetUserKeys Get the keys for all users<获取所有的key>
func (get *Client) GetUserKeys() (userKeys []string) {
	userKeys = make([]string, 0)
	get.UserLock.RLock()
	for key := range get.Users {
		userKeys = append(userKeys, key)
	}
	defer get.UserLock.RUnlock()
	return
}

//GetUserList 获取用户的key
func (get *Client) GetUserList(appId string) (userList []string) {
	userList = make([]string, 0)
	get.UserLock.RLock()
	defer get.UserLock.RUnlock()
	for _, v := range get.Users {
		if v.AppId == appId {
			userList = append(userList, v.UserId)
		}
	}
	return
}

//GetUserClients 获取用户的连接
func (get *Client) GetUserClients() (connection []*Connection) {

	connection = make([]*Connection, 0)
	get.UserLock.RLock()

	for _, v := range get.Users {
		connection = append(connection, v)
	}
	defer get.UserLock.RUnlock()
	return
}

//SendAll 向全部成员(除了自己)发送数据
func (get *Client) SendAll(message []byte, connection *Connection) {
	clients := get.GetUserClients()
	for _, conn := range clients {
		if conn != connection {
			conn.WriteMessage(message)
		}
	}
}

//OnRegister 用户建立连接事件
func (get *Client) OnRegister(connection *Connection) {
	get.AddClients(connection)
}

//OnLogin 用户登录
func (get *Client) OnLogin(login *Login) {
	var client = login.Client
	// 连接存在，在添加
	if get.IsClient(client) {
		userKey := login.GetUserKey()
		get.AddUsers(userKey, login.Client)
	}
}

//OnUserLogout 用户断开连接
func (get *Client) OnUserLogout(client *Connection) {
	get.DeleteClients(client)
	deleteResult := get.DeleteUsers(client)
	if deleteResult == false {
		return
	}
}
