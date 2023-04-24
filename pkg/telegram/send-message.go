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

var activeChats = make(map[int64]bool)

func SendMessage(text string, bot *tgbotapi.BotAPI) {
	BotToken := os.Getenv("BOT_TOKEN")
	TelegramURL := viper.GetString("telegramURL")

	// get latest update
	updates, err := bot.GetUpdates(tgbotapi.NewUpdate(0))
	if err != nil {
		log.Panic(err)
	}

	var chatIds int64
	if len(updates) > 0 {
		chatIds = updates[len(updates)-1].Message.Chat.ID
	}

	if !activeChats[chatIds] {
		activeChats[chatIds] = true
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
