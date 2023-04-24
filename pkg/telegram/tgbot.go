package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config struct {
	BotToken    string
	TelegramUrl string
}

func NewTgBot(cfg Config) (*tgbotapi.BotAPI, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	return bot, nil
}
