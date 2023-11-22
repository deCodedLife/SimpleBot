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
	s.handlers = make(map[string]MessageHandler)
	s.ApiUrl = u
	s.Token = t
	return s
}

func (s *Bot) Contains(m model.Message) bool {
	for _, poolMessage := range s.pool {
		if poolMessage.MessageId == m.MessageId {
			return true
		}
	}
	return false
}

func (s *Bot) AddHandler(m string, h MessageHandler) {
	s.handlers[m] = h
}

func (s *Bot) Add(m model.Message) {
	if s.handlers[strings.ToLower(m.Text)] != nil {
		s.handlers[strings.ToLower(m.Text)](m.Chat)
	}
	s.pool = append(s.pool, m)
}

func (s *Bot) SendMessage(c model.Chat, m string) {
	var response ReplyMessage
	response.ChatId = strconv.Itoa(c.Id)
	response.Message = m

	buff, err := json.Marshal(response)

	if err != nil {
		log.Println("Не удалось сформировать сообщение по причине: ", err.Error())
		return
	}

	_, err = http.Post(fmt.Sprintf("%s%s/sendMessage", s.ApiUrl, s.Token), "application/json", bytes.NewBuffer(buff))

	if err != nil {
		println("Не удалось отправить сообщение по причине: ", err.Error())
	}
}

func (s *Bot) LoadMessages(m []model.Message) {
	s.pool = m
}

func (s *Bot) Poll(c chan model.Message, m string, d time.Duration) {

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
				c <- msgCtx.Message
			}

			time.Sleep(time.Second * d)

		}

	}()

}
