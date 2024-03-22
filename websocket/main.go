package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// 所有在线客户端
	clients = make(map[string]*websocket.Conn)

	// 设置客户端超时时间
	timeoutDuration = 1 * time.Minute
)

func main() {
	engine := gin.Default()

	engine.GET("/hello", func(ctx *gin.Context) {
		token := ctx.GetHeader("token")

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			fmt.Println("Error upgrading to websocket:", err)
			return
		}
		clients[token] = conn
		defer func() {
			conn.Close()
			delete(clients, token)
		}()

		for {
			select {
			case <-time.After(timeoutDuration):
				log.Println(token, "web socket connection timeout")
				return
			default:
				conn.SetReadDeadline(time.Now().Add(timeoutDuration))
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println(token, "socket 断开了:", err)
					return
				}
				fmt.Println("收到消息", string(msg))
				for t, c := range clients {
					if t == token {
						continue
					}
					fmt.Println("发送消息给", t)
					c.WriteMessage(websocket.TextMessage, msg)
				}
			}
		}
	})

	engine.Run(":8080")
}
