package main

import (
	"bot_lab/internal/bot"
	"bot_lab/internal/db"
	"bot_lab/internal/db/local"
	"bot_lab/internal/model"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	API_URL   = "https://api.telegram.org/bot"
	BOT_TOKEN = "6739808790:AAEmT1ZgyyBkRJzd1jg-FA-16s8NjTzHmg0"

	Bot         bot.Bot
	Database    db.DatabaseHandler
	DB_SETTINGS db.Credentials
)

func initHandlers() {
	Bot.AddHandler(strings.ToLower("Hello"), func(m model.Message) {
		Bot.SendMessage(m.Chat, "You too")
	})
	Bot.AddHandler(strings.ToLower("Привет"), func(m model.Message) {
		message := fmt.Sprintf("Приветствуем Вас %s", m.From.FirstName)
		Bot.SendMessage(m.Chat, message)
	})
	Bot.AddHandler(strings.ToLower("Сколько сейчас времени?"), func(m model.Message) {
		currentTime := fmt.Sprintf("%d-%d-%d %d:%d", time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute())
		Bot.SendMessage(m.Chat, currentTime)
	})
}

func main() {
	log.Println("[ ] Simple bot starting ")

	DB_SETTINGS.FileName = "localDB.json"

	log.Print("[ ] Инициализация базы данных ")
	Database = &local.Database{}
	err := Database.NewConnection(DB_SETTINGS)

	if err != nil {
		log.Panicln("Невозможно подключиться к базе данных", err)
	}

	log.Println("[ ] Чтение базы данных ")
	rows := Database.GetMessages()

	Bot.New(API_URL, BOT_TOKEN)
	Bot.LoadMessages(rows)
	initHandlers()

	log.Println("[ ] Прослушивание входящих сообщений")

	messages := make(chan model.Message)
	Bot.Poll(messages, "getUpdates", time.Second)

	for {
		message := <-messages
		if Bot.Contains(message) {
			continue
		}

		Bot.HandleMessage(message)
		_ = Database.SaveMessage(message)

		log.Println(fmt.Sprintf("Получено сообщение от пользователя %s: %s", message.From.Username, message.Text))
	}

}
