package bot

import (
	"bot_lab/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Bot) New(u string, t string) *Bot {
	s.messageHandlers = make(map[string]MessageHandler)
	s.callbackHandlers = make(map[string]ContextHandler)
	s.ApiUrl = u
	s.Token = t
	return s
}

func (s *Bot) Contains(m model.MessageContext) bool {
	for _, poolMessage := range s.pool {
		if poolMessage.UpdateId == m.UpdateId {
			return true
		}
	}
	return false
}

func (s *Bot) AddMessageHandler(m string, h MessageHandler) {
	s.messageHandlers[m] = h
}

func (s *Bot) AddCallbackHandler(m string, h ContextHandler) {
	s.callbackHandlers[m] = h
}

func (s *Bot) Handle(m model.MessageContext) {

	messageText, _ := strconv.Unquote("\"" + m.Message.Text + "\"")
	messageText = strings.ToLower(messageText)

	if s.messageHandlers[messageText] != nil {
		s.messageHandlers[messageText](m.Message)
	}

	messageText, _ = strconv.Unquote("\"" + m.CallbackQuery.Data + "\"")
	messageText = strings.ToLower(messageText)

	if s.callbackHandlers[messageText] != nil {
		s.callbackHandlers[messageText](m.CallbackQuery)
	}

	s.pool = append(s.pool, m)

}

func (s *Bot) SendMessage(m ReplyMessage) {
	buff, err := json.Marshal(m)

	if err != nil {
		log.Println("Не удалось сформировать сообщение по причине: ", err.Error())
		return
	}

	reply, err := http.Post(fmt.Sprintf("%s%s/sendMessage", s.ApiUrl, s.Token), "application/json", bytes.NewBuffer(buff))

	if err != nil {
		println("Не удалось отправить сообщение по причине: ", err.Error())
	}

	var response interface{}
	_ = json.NewDecoder(reply.Body).Decode(&response)

}

func (s *Bot) LoadMessages(m []model.MessageContext) {
	s.pool = m
}

func (s *Bot) Poll(c chan model.MessageContext, m string, d time.Duration) {

	go func() {

		for {

			url := fmt.Sprintf("%s%s/%s", s.ApiUrl, s.Token, m)
			resp, err := http.Get(url)

			if err != nil {
				log.Println("При запросе данных произошла ошибка: ", err.Error())
			}

			var response Response
			err = json.NewDecoder(resp.Body).Decode(&response)

			if err != nil {
				log.Println("Получена неверная структура данных: ", err.Error())
			}

			for _, msgCtx := range response.Result {
				c <- msgCtx
			}

			time.Sleep(d)

		}

	}()

}
