package telegram

import (
	"ext-tg-bot/utils"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func (bot *TgBot) Log(incoming *tgbotapi.Message) {
	logMsg := fmt.Sprintf("<b>New message from: %d, @%s</b>\n%s", incoming.From.ID, incoming.From.UserName, incoming.Text)
	fmt.Println(bot.logID)
	msg := tgbotapi.NewMessage(bot.logID, logMsg)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.api.Send(msg)
	if err != nil {
		utils.LogWarn("message from:", incoming.From.ID, err)
	}
}
