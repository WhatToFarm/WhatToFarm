package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"ton-tg-bot/logger"
	"ton-tg-bot/models"
)

type TgBot struct {
	api   *tgbotapi.BotAPI
	users map[int64]*models.TgUser
	logID int64
}

func (bot *TgBot) Init(botID string, logID int64) {
	bot.init(botID, logID)
}

func (bot *TgBot) init(botID string, logID int64) {
	if err := bot.getUsers(); err != nil {
		logger.LogFatal("Users from DB:", err)
	}

	api, err := tgbotapi.NewBotAPI(botID)
	if err != nil {
		logger.LogFatal("Bot API:", err)
	}
	bot.api = api
	bot.api.Debug = false
	bot.logID = logID

	logger.LogInfo("Authorized on account:", bot.api.Self.UserName)
}

func (bot *TgBot) Start() {
	bot.start()
}

func (bot *TgBot) start() {
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
