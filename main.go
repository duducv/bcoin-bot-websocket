package main

import (
	"crypto-bot-websocket/domain"
	"crypto-bot-websocket/repository"
	"crypto-bot-websocket/usecase"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var repo = repository.NewClientRepository()

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		panic("canno't generate conn")
	}

	authMessageType, authPayload, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	err = usecase.CheckAuthorization(authPayload)

	if err != nil {
		conn.WriteMessage(authMessageType, []byte(err.Error()))
		conn.Close()
	}

	conn.WriteMessage(authMessageType, []byte("access granted"))

	handleMessageUsecase := usecase.HandleMessageUsecase{Repository: repo}
	newClient := domain.NewClient(conn)

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			handleMessageUsecase.RemoveClient(newClient)
			disconnectMessage := fmt.Sprintf("disconnected, total clients: %d", len(*handleMessageUsecase.Repository.GetAllClients()))
			fmt.Println(disconnectMessage)
			break
		}

		conn.WriteMessage(messageType, []byte(payload))
		err = handleMessageUsecase.Process(messageType, payload, newClient)
		if err != nil {
			conn.WriteMessage(messageType, []byte(err.Error()))
			conn.Close()
		}
	}

}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("error at load .env file")
	}

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":5555", nil)
}
