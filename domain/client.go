package domain

import (
	"github.com/gorilla/websocket"
	"github.com/satori/uuid"
)

type Client struct {
	ConnectionUniqueID string
	ID                 int
	Conn               *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		uuid.NewV4().String(),
		0,
		conn,
	}
}

func (c *Client) SetId(id int) {
	c.ID = id
}

type IClientRepository interface {
	AddClient(client *Client)
	RemoveClient(client *Client)
	GetAllClients() *[]Client
}
