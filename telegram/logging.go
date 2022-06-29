package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"ton-tg-bot/logger"
)

func (bot *TgBot) forwardMessage(message *tgbotapi.Message) {
	logMsg := fmt.Sprintf("<b>New message from: %d, @%s</b>\n%s", message.From.ID, message.From.UserName, message.Text)
	fmt.Println(bot.logID)
	msg := tgbotapi.NewMessage(bot.logID, logMsg)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.api.Send(msg)
	if err != nil {
		logger.LogWarn("message to:", message.From.ID, err)
	}
}

func (bot *TgBot) logToGroup(message *tgbotapi.Message, text string) {
	logMsg := fmt.Sprintf("<b>Log. User: %d, @%s</b>\n%s", message.From.ID, message.From.UserName, text)
	fmt.Println(bot.logID)
	msg := tgbotapi.NewMessage(bot.logID, logMsg)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.api.Send(msg)
	if err != nil {
		logger.LogWarn("message to:", message.From.ID, err)
	}
}
