package telegram

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

func SendMessage(text string) {
	// init bot
	BotToken := os.Getenv("BOT_TOKEN")
	TelegramURL := viper.GetString("telegramURL")

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	var chatId int64 // get user id

	// find out chat id and send info about bot
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message != nil { // If we got a message
			if update.Message.Text == "/start" || update.Message.Text == "/info" {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				chatId = update.Message.Chat.ID

				textMess := "Hello! Thank you for joining us!\n\n"
				textMess += "At the moment, vacancies are sent only from Indeed.com, DOU.ua\n\n"
				textMess += "We plan to add such services in the future:\n"
				textMess += " - Robota.ua\n - Inco.works\n - Work.ua"

				msg := tgbotapi.NewMessage(chatId, textMess)

				bot.Send(msg)
			}
		}
	}

	// send message
	textConv := fmt.Sprintf("%s", text)
	textJson := fmt.Sprintf(`{"chat_id":%d, "text":"%s", "parse_mode":"HTML", "disable_web_page_preview": true}`, chatId, textConv)
	data := []byte(textJson)

	txt := bytes.NewReader(data)
	req := fmt.Sprintf("%s%s/sendMessage", TelegramURL, BotToken)

	_, err = http.Post(req, "application/json", txt)
	if err != nil {
		log.Panic(err)
	}
	// logging
	// ...
}
