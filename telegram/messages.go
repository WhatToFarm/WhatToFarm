package telegram

import (
	"ext-tg-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *TgBot) sendMessage(ID int64, message func() string) {
	msg := tgbotapi.NewMessage(ID, message())
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	_, err := bot.api.Send(msg)
	if err != nil {
		utils.LogWarn(ID, err)
	}
}

func wrong() string {
	return fmt.Sprintf("Something wrong, try again")
}
