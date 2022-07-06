package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"ton-tg-bot/github"
	"ton-tg-bot/logger"
	"ton-tg-bot/models"
	"ton-tg-bot/mongo"
)

const (
	StateInitiated = iota
	StateValidation
	StateValidateGitHub

	Start = "/start"
)

func (bot *TgBot) extractData(message *tgbotapi.Message) {
	if message.From.IsBot {
		bot.sendMessage(message.Chat.ID, isBot)
		return
	}

	user := bot.users[message.Chat.ID]
	if user == nil {
		user = &models.TgUser{
			TgId:  message.Chat.ID,
			State: StateValidation,
		}
		bot.users[message.Chat.ID] = user
	}

	bot.checkState(message, user)
}

func (bot *TgBot) checkState(message *tgbotapi.Message, user *models.TgUser) {
	switch user.State {
	case StateValidation:
		bot.sendMessage(user.TgId, func() string { return validationStepOne(user.TgId) })
		user.State = StateValidateGitHub
	case StateValidateGitHub:
		bot.parseGitHubName(message, user)
	default:
		bot.parseMessage(message, user)
	}
}

func (bot *TgBot) parseGitHubName(message *tgbotapi.Message, user *models.TgUser) {
	user.GitAccount = strings.TrimSpace(message.Text)
	if err := github.ValidateGitHub(user.GitAccount, user.TgId); err != nil {
		logger.LogWarn("GitHub validation user:", user.TgId, "error:", err)
		bot.sendMessage(user.TgId, func() string {
			return validationFailed(err.Error())
		})
		user.State = StateValidation
		return
	}
	if err := mongo.CreateUser(user); err != nil {
		logger.LogWarn("GitHub validation user:", user.TgId, "error:", err)
		bot.sendMessage(user.TgId, func() string {
			return validationFailed(err.Error())
		})
		user.State = StateValidation
		return
	}

	user.State = StateInitiated
	bot.sendMessage(user.TgId, validationSuccess)
	bot.logText(message, fmt.Sprintf("GitHub user <code>github.com/%s</code> registered", message.Text))
}

func (bot *TgBot) parseMessage(message *tgbotapi.Message, user *models.TgUser) {
	start := strings.Contains(message.Text, Start)
	switch {
	case start:
		bot.sendMessage(user.TgId, description)
	case message.Document != nil:
		bot.handleFileUpload(message, user)
	case message.Text != "":
		bot.logForward(message)
		bot.sendMessage(user.TgId, answer)
	}
}

func (bot *TgBot) getUsers() error {
	bot.users = make(map[int64]*models.TgUser)
	users, err := mongo.GetUsers()
	if err != nil {
		return err
	}
	for i, u := range users {
		bot.users[u.TgId] = users[i]
	}
	return nil
}
