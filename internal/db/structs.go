package db

import "bot_lab/internal/model"

type Credentials struct {
	FileName string
	Host     string
	Port     string
	Username string
	Password string
}

type DatabaseHandler interface {
	NewConnection(c Credentials) error
	DropConnection()
	SaveMessage(m model.MessageContext) error
	GetMessages() []model.MessageContext
}
