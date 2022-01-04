package usecase

import (
	"crypto-bot-websocket/domain"
	"crypto-bot-websocket/dto"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var CHECK = "check"
var REMOVE = "remove"

type HandleMessageUsecase struct {
	Repository domain.IClientRepository
}

func (usecase HandleMessageUsecase) Process(messageType int, payload []byte, client *domain.Client) error {
	jsonPayload := dto.NewMessageDto()
	if err := usecase.matchJsonPayload(payload, &jsonPayload); err != nil {
		return err
	}

	client.SetId(jsonPayload.ID)

	if !usecase.checkIfAlreadyExists(client, messageType) {
		log.Println(fmt.Sprintln("clients %i", len(*usecase.Repository.GetAllClients())))
	}

	usecase.AddClient(client)

	connectedMessage := fmt.Sprintf("conectado, total de clients: %d", len(*usecase.Repository.GetAllClients()))
	fmt.Println(connectedMessage)

	return nil
}

func (usecase HandleMessageUsecase) checkIfAlreadyExists(client *domain.Client, messageType int) bool {
	for _, clientRange := range *usecase.Repository.GetAllClients() {
		if client.ID == clientRange.ID {
			clientRange.Conn.WriteMessage(messageType, []byte("found simultaneous connections for this id"))
			client.Conn.WriteMessage(messageType, []byte("found simultaneous connections for this id"))
			return true
		}
	}
	client.Conn.WriteMessage(messageType, []byte("connected with no restrictions"))
	return false
}

func (usecase HandleMessageUsecase) RemoveClient(client *domain.Client) {
	usecase.Repository.RemoveClient(client)
}

func (usecase HandleMessageUsecase) AddClient(client *domain.Client) {
	usecase.Repository.AddClient(client)
}

func (usecase HandleMessageUsecase) matchJsonPayload(payload []byte, dto *dto.MessageDto) error {

	err := json.Unmarshal(payload, &dto)

	if err != nil {
		return errors.New("failed to decode json")
	}

	return nil
}

func CheckAuthorization(payload []byte) error {
	keyDto := dto.NewKeyDto()
	err := json.Unmarshal(payload, &keyDto)

	if err != nil {
		return errors.New("failed to decode key json")
	}

	key := os.Getenv("WS_KEY")
	if key != keyDto.Key {
		return errors.New("key invalid")
	}

	return nil
}
