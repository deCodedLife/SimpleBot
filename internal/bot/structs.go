package bot

import "bot_lab/internal/model"

type MessageContext struct {
	UpdateId int           `json:"update_id"`
	Message  model.Message `json:"message"`
}

type Response struct {
	Ok     bool             `json:"ok"`
	Result []MessageContext `json:"result"`
}

type MessageHandler func(m model.Message)

type ReplyMessage struct {
	ChatId  string `json:"chat_id"`
	Message string `json:"text"`
}

type Bot struct {
	ApiUrl string
	Token  string

	pool     []model.Message
	handlers map[string]MessageHandler
}
