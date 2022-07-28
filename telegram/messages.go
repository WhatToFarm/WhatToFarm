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
	return "Sorry. Internal error, try again."
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

func manyAttempts(minutes int64) string {
	return fmt.Sprintf("Sorry, you have only 5 attempts per hour.\n"+
		"Next time you can try in <b>%d minutes</b>.", minutes)
}

func bigFileSize() string {
	return "To big file size (expected 1MB maximum)"
}

func validationStepOne(id int64) string {
	return fmt.Sprintf("Your GitHub account should be created a month ago minimum.\n"+
		"It's important!\n\n"+
		"Create public repository named <b>\"TCS2-%d\"</b> "+
		"in your GitHub. It may be empty.\n"+
		"So, send me your GitHub account name (not the link to repository) after this message.\n"+
		"I will let you know next steps!.", id)
}

func validationFailed(err string) string {
	return fmt.Sprintf("I can't register your account.\n"+
		"Error: %v\n"+
		"Please, try again", err)
}

func validationSuccess() string {
	return "Everything is OK!\n" +
		"You can send me your solution or questions."
}

func description() string {
	return "You have validated your GitHug account and can continue.\n" +
		"You can send me your solution or questions.\n\n" +
		"Your solution must include only <b>tar.gz</b> archive file.\n" +
		"If you have any questions, message me. And specialist from our teem contact you."
}

func waitResult() string {
	return "Started solution evaluation."
}

func expireDeadline() string {
	return "Solution acceptance is closed"
}
