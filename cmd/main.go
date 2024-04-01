package main

import (
	"bot_lab/internal/bot"
	"bot_lab/internal/db"
	"bot_lab/internal/db/local"
	"bot_lab/internal/model"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	API_URL   = "https://api.telegram.org/bot"
	BOT_TOKEN = "7154811283:AAGb4a4eUfpHY-e4PkRXIARJ_EL-wAyn5jI"

	Bot         bot.Bot
	Database    db.DatabaseHandler
	DB_SETTINGS db.Credentials
)

func initHandlers() {
	Bot.AddMessageHandler(strings.ToLower("/start"), func(m model.Message) {
		var response bot.ReplyMessage
		response.ChatId = strconv.Itoa(m.Chat.Id)
		response.Message = "You too"

		var buttons []bot.InlineKeyboardButton
		buttons = append(buttons, bot.InlineKeyboardButton{
			Text:         "Подтвердить",
			CallbackData: "confirm",
		})
		buttons = append(buttons, bot.InlineKeyboardButton{
			Text:         "Отменить",
			CallbackData: "cancel",
		})
		response.ReplyMarkup.InlineKeyboard = append(response.ReplyMarkup.InlineKeyboard, buttons)

		Bot.SendMessage(response)
	})
	Bot.AddMessageHandler(strings.ToLower("Привет"), func(m model.Message) {
		var response bot.ReplyMessage
		response.ChatId = strconv.Itoa(m.Chat.Id)
		response.Message = fmt.Sprintf("Приветствуем Вас %s", m.From.FirstName)
		Bot.SendMessage(response)
	})
	Bot.AddMessageHandler(strings.ToLower("Сколько сейчас времени?"), func(m model.Message) {
		var response bot.ReplyMessage
		response.ChatId = strconv.Itoa(m.Chat.Id)
		response.Message = fmt.Sprintf("%d-%d-%d %d:%d", time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute())
		Bot.SendMessage(response)
	})
	Bot.AddCallbackHandler("confirm", func(m model.CallbackQuery) {
		var response bot.ReplyMessage
		response.ChatId = strconv.Itoa(m.Message.Chat.Id)
		response.Message = fmt.Sprintf("Ждём Вас на приёме %s", m.From.FirstName)
		Bot.SendMessage(response)
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

	messages := make(chan model.MessageContext)
	Bot.Poll(messages, "getUpdates", time.Second)

	for {
		message := <-messages

		if Bot.Contains(message) {
			continue
		}

		Bot.Handle(message)
		_ = Database.SaveMessage(message)

		log.Println(fmt.Sprintf("Получено сообщение от пользователя %s: %s", message.Message.From.Username, message.Message.Text))
	}

	//_ = http.ListenAndServe(":8080", nil)

}
