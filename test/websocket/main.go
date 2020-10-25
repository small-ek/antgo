package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/small-ek/ginp/net/gwebsocket"
	"log"
	"net/http"
	"time"
)

var upGrader = (&websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
})

func main() {
	bindAddress := "127.0.0.1:1111"
	r := gin.Default()
	r.GET("/ping", ping)
	r.Run(bindAddress)
}

//webSocket请求ping 返回pong
func ping(c *gin.Context) {
	var (
		websocket *websocket.Conn
		err       error
		conn      *gwebsocket.Connection
		data      []byte
	)
	// 完成ws协议的握手操作
	if websocket, err = upGrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		return
	}
	conn = gwebsocket.New(websocket, c.ClientIP(), uint64(time.Now().Unix()))
	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
		log.Println(string(data))
	}
	gwebsocket.NewClient().Register <- conn

ERR:
	conn.Close()
}
