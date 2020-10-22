package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/small-ek/ginp/net/gwebsocket"
	"log"
	"net/http"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	bindAddress := "127.0.0.1:1234"
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

	if conn, err = gwebsocket.New(websocket); err != nil {
		goto ERR
	}
	// 启动线程，不断发消息
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
		log.Println(string(data))
	}

ERR:
	conn.Close()
}
