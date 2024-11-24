package main

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"log"
	"net/http"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {

	socketServer := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	socketServer.OnConnect("/", func(s socketio.Conn) error {
		log.Println("Connected:", s.ID())
		s.SetContext("")
		s.Join("chat")
		return nil
	})

	socketServer.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Println("Message received:", msg)
		socketServer.BroadcastToRoom("/", "chat", "message", msg)
	})

	socketServer.OnEvent("/", "chat", func(s socketio.Conn, reason string) {
		log.Println("Disconnect:", s.ID(), reason)
	})

	go socketServer.Serve()
	defer socketServer.Close()

	router := gin.Default()

	router.GET("/socket.io/*any", gin.WrapH(socketServer))
	router.POST("/socket.io/*any", gin.WrapH(socketServer))

	router.Run(":4000")
}
