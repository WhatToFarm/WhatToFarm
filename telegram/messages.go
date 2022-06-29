package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"ton-tg-bot/logger"
)

func (bot *TgBot) sendMessage(ID int64, message func() string) {
	msg := tgbotapi.NewMessage(ID, message())
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	_, err := bot.api.Send(msg)
	if err != nil {
		logger.LogWarn(ID, err)
	}
}

func wrong() string {
	return "Internal error, try again."
}

func invalidFormat() string {
	return fmt.Sprintf("Invalid file format, expected \"%s\"", defaultExt)
}

func answer() string {
	return "We have received your message and will contact you as soon as possible."
}

func isBot() string {
	return "I see you are bot."
}

func toManyAttempts(minutes int64) string {
	return fmt.Sprintf("Sorry, you have only 5 attempts per hour."+
		"Next time you can try in <b>%d minutes</b>.", minutes)
}

func validateUser() string {
	return "You need authorize your GitHub account in bot."
}

func validateStepOne() string {
	return "Your GitHub account should be created a month ago minimum.\n" +
		"It's important!\n\n" +
		"Create public repository named <b>\"TCS2\"</b> in your GitHub. It may be empty.\n" +
		"So, send me your GitHub account name after this message.\n" +
		"If everything is OK, I will let you know!."
}

func validateFailed(err string) string {
	return fmt.Sprintf("I can't register your account.\n"+
		"Error: %s\n"+
		"Please, try again", err)
}
