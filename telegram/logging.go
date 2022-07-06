package telegram

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"ton-tg-bot/logger"
)

func (bot *TgBot) logText(message *tgbotapi.Message, text string) {
	logMsg := fmt.Sprintf("<b>Log. User: %d, @%s</b>\n%s", message.From.ID, message.From.UserName, text)
	msg := tgbotapi.NewMessage(bot.logID, logMsg)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.api.Send(msg)
	if err != nil {
		logger.LogWarn("message to log chat:", err)
	}
}

func (bot *TgBot) logForward(message *tgbotapi.Message) {
	msg := tgbotapi.NewForward(bot.logID, message.Chat.ID, message.MessageID)
	_, err := bot.api.Send(msg)
	if err != nil {
		logger.LogWarn("message to log chat:", err)
	}
}
