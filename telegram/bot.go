package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"ton-tg-bot/logger"
	"ton-tg-bot/models"
	"ton-tg-bot/mongo"
)

func (bot *TgBot) extractData(message *tgbotapi.Message) {
	if message.From.IsBot {
		bot.sendMessage(message.Chat.ID, isBot)
		return
	}

	logger.LogDebug("extractData: ", message.From.UserName, message.Text)
	if message.Text != "" {
		bot.ForwardMessage(message)
	}

	userData := bot.users[message.Chat.ID]
	if userData == nil {
		userData = &models.TgUser{
			TgId:  message.Chat.ID,
			State: StateValidation,
		}
		bot.users[message.Chat.ID] = userData
	}

	bot.checkState(message, userData)
}

func (bot *TgBot) checkState(message *tgbotapi.Message, userData *models.TgUser) {
	switch userData.State {
	case StateValidation:
		bot.sendMessage(userData.TgId, validateUser)
		bot.sendMessage(userData.TgId, validateStepOne)
		userData.State = StateValidateGitHub
	case StateValidateGitHub:
		bot.parseGitHubName(message, userData)
	default:
		bot.parseMessage(message, userData)
	}
}

func (bot *TgBot) parseGitHubName(message *tgbotapi.Message, userData *models.TgUser) {
	userData.GitAccount = strings.TrimSpace(message.Text)
	userData.State = StateInitiated
	if err := mongo.CreateUser(userData); err != nil {
		bot.sendMessage(userData.TgId, func() string {
			return validateFailed(err.Error())
		})
		userData.State = StateValidation
	}
}

func (bot *TgBot) parseMessage(message *tgbotapi.Message, userData *models.TgUser) {
	switch {
	case message.Text != "":
		bot.sendMessage(userData.TgId, answer)
	case message.Document != nil:
		bot.handleFileUpload(message, userData.TgId)
	}
}

func (bot *TgBot) getUsers() error {
	users, err := mongo.GetUsers()
	if err != nil {
		return err
	}
	for i, u := range users {
		bot.users[u.TgId] = users[i]
	}
	return nil
}
