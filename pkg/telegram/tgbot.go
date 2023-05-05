package telegram

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"net/http"
	tg_methods "vac_informer_tgbot/pkg/database/telegram_db_methods"
)

type Config struct {
	BotToken    string
	TelegramUrl string
}

var (
	bot         *tgbotapi.BotAPI
	BotToken    string
	activeChats = make(map[int64]bool)
)

func NewTgBot(cfg Config) (*tgbotapi.BotAPI, error) {
	var err error

	BotToken = cfg.BotToken

	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	return bot, nil
}

func SendMessage(text string) {
	TelegramURL := viper.GetString("telegramURL")

	// get latest update
	updates, err := bot.GetUpdates(tgbotapi.NewUpdate(0))
	if err != nil {
		log.Panic(err)
	}

	var setChatIds int64
	if len(updates) > 0 {
		setChatIds = updates[len(updates)-1].Message.Chat.ID
	}
	tg_methods.CheckUser(setChatIds)

	chatIds, err := tg_methods.SelectAllUsers()
	for _, chatId := range chatIds {
		activeChats[chatId] = true
	}

	// send message
	for chatId := range activeChats {
		textConv := fmt.Sprintf("%s", text)
		textJson := fmt.Sprintf(`{"chat_id":%d, "text":"%s", "parse_mode":"HTML", "disable_web_page_preview": false}`, chatId, textConv)
		data := []byte(textJson)

		txt := bytes.NewReader(data)
		req := fmt.Sprintf("%s%s/sendMessage", TelegramURL, BotToken)

		_, err = http.Post(req, "application/json", txt)
		if err != nil {
			log.Panic(err)
		}
	}
}
