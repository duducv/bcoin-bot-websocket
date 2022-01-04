package repository

import (
	"crypto-bot-websocket/domain"
)

type ClientRepository struct {
	Clients []domain.Client
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (repository *ClientRepository) AddClient(client *domain.Client) {
	if !repository.checkIfActualConnUniqueIDexists(client) {
		repository.Clients = append(repository.Clients, *client)
	}
}

func (repository *ClientRepository) checkIfActualConnUniqueIDexists(client *domain.Client) bool {
	for _, actualClient := range repository.Clients {
		if actualClient.ConnectionUniqueID == client.ConnectionUniqueID {
			return true
		}
	}
	return false
}

func (repository *ClientRepository) RemoveClient(client *domain.Client) {
	for index, actualClient := range repository.Clients {
		if actualClient.ConnectionUniqueID == client.ConnectionUniqueID {
			repository.Clients = append(repository.Clients[:index], repository.Clients[index+1:]...)
		}
	}
}

func (repository *ClientRepository) GetAllClients() *[]domain.Client {
	return &repository.Clients
}
