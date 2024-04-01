package bot

import "bot_lab/internal/model"

type Response struct {
	Ok     bool                   `json:"ok"`
	Result []model.MessageContext `json:"result"`
}

type MessageHandler func(m model.Message)
type ContextHandler func(m model.CallbackQuery)

type ReplyMessage struct {
	ChatId      string               `json:"chat_id"`
	Message     string               `json:"text"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	Url          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data"`
}

type Bot struct {
	ApiUrl string
	Token  string

	pool             []model.MessageContext
	messageHandlers  map[string]MessageHandler
	callbackHandlers map[string]ContextHandler
}
