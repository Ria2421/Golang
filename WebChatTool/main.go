package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Type    int
	Message []byte
}

func main() {
	r := gin.Default()

	// websocketのupgraderを定期
	wsupgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// TOPページ
	r.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "index.html")
	})

	r.GET("/ws", func(ctx *gin.Context) {
		// upgraderを呼び出すことで通常のhttp通信からwebsocketへupgrade
		// コネクションを作成する
		conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %+v\n", err)
			return
		}

		// コネクションをclientsマップへ追加
		clients[conn] = true

		// 無限ループさせることでクライアントからのメッセージを受け付けられる状態にする
		// クライアントとのコネクションが切れた場合はReadMessage()関数からエラーが返る
		for {
			t, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("ReadMessage Error. ERROR: %+v\n", err)
				break
			}
			// 受け取ったメッセージをbroadcastを通じてhandleMessages()関数へ渡す
			broadcast <- Message{Type: t, Message: msg}
		}
	})

	// 非同期でhandleMessagesを実行
	go handleMessages()

	r.Run(":4001")
}

// broadcastにメッセージがあれば、clientsに格納されている全てのコネクションへ送信する
func handleMessages() {
	for {
		message := <-broadcast

		for client := range clients {
			err := client.WriteMessage(message.Type, message.Message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
