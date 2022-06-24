package telegram

import (
	"ext-tg-bot/models"
	"ext-tg-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TgBot struct {
	api    *tgbotapi.BotAPI
	states map[int64]*models.TgUser
	logID  int64
}

func (bot *TgBot) Init(botID string, logID int64) {
	bot.init(botID, logID)
}

func (bot *TgBot) init(botID string, logID int64) {
	bot.states = make(map[int64]*models.TgUser)
	api, err := tgbotapi.NewBotAPI(botID)
	if err != nil {
		utils.LogFatal(err)
	}
	bot.api = api
	bot.api.Debug = false
	bot.logID = logID

	utils.LogInfo("Authorized on account:", bot.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	_, _ = bot.api.RemoveWebhook()
	updates, _ := bot.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		go bot.extractData(update.Message)
	}
}

func (bot *TgBot) extractData(message *tgbotapi.Message) {
	utils.LogDebug("extractData: ", message.From.UserName, message.Text)
	if message.Text != "" {
		bot.Log(message)
	}

	userData := bot.states[message.Chat.ID]
	if userData == nil {
		userData = &models.TgUser{
			TgId:  message.Chat.ID,
			State: models.StateInitial,
		}
		bot.states[message.Chat.ID] = userData
	}

	bot.parseMessage(message, userData)
}

func (bot *TgBot) parseMessage(message *tgbotapi.Message, userData *models.TgUser) {
	switch {
	case message.Text != "":
		userData.State = models.StateMessage
	case message.Document != nil:
		bot.handleFileUpload(message, userData.TgId)
		userData.State = models.StateInitial
	}
	bot.reply(message, userData)
}

func (bot *TgBot) reply(incoming *tgbotapi.Message, userData *models.TgUser) {
	switch userData.State {
	case models.StateInitial:
		//bot.sendMessage(incoming.Chat.ID, func() string {
		//	return "11111"
		//})
	case models.StateMessage:
		//bot.sendMessage(incoming.Chat.ID, func() string {
		//	return "22222"
		//})
	}
}
