package main

import (
	"bot_lab/internal/bot"
	"bot_lab/internal/db"
	"bot_lab/internal/db/local"
	"bot_lab/internal/model"
	"fmt"
	"log"
)

var (
	API_URL     = "https://api.telegram.org/bot"
	BOT_TOKEN   = "6739808790:AAEmT1ZgyyBkRJzd1jg-FA-16s8NjTzHmg0"
	Bot         bot.Bot
	Database    db.DatabaseHandler
	DB_SETTINGS db.Credentials
)

func Reply(m model.Chat) {
	Bot.SendMessage(m, "You too")
}

func main() {
	log.Println("Started")

	DB_SETTINGS.FileName = "localDB.json"

	Database = &local.Database{}
	err := Database.NewConnection(DB_SETTINGS)
	rows := Database.GetMessages()

	if err != nil {
		log.Panicln("Невозможно подключиться к базе данных", err)
	}

	Bot.New(API_URL, BOT_TOKEN)
	Bot.LoadMessages(rows)
	Bot.AddHandler("hello", Reply)

	messages := make(chan model.Message)
	Bot.Poll(messages, "getUpdates", 3)

	for {
		message := <-messages
		if Bot.Contains(message) {
			continue
		}
		Bot.Add(message)
		_ = Database.SaveMessage(message)

		log.Println(fmt.Sprintf("Получено сообщение от пользователя %s: %s", message.From.Username, message.Text))
	}

}
