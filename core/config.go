package core

import (
	"github.com/spf13/viper"

	"ton-tg-bot/logger"
)

var (
	MongoURL      string
	MongoDatabase string
	BotId         string
	LogId         int64
	Host          string
	BasePath      string
)

func Init() {
	MongoURL = viper.GetString("mongo.url")
	MongoDatabase = viper.GetString("mongo.database")
	Host = viper.GetString("service.host")
	BasePath = viper.GetString("service.base_path")

	BotId = viper.GetString("telegram.bot_id")
	if BotId == "" {
		logger.LogFatal("Telegram bot ID not provided")
	}

	LogId = viper.GetInt64("telegram.log_id")
	if LogId == 0 {
		logger.LogFatal("Telegram log ID not provided")
	}

	logger.LogInfo("CFG: MongoURL = ", MongoURL)
	logger.LogInfo("CFG: MongoDatabase = ", MongoDatabase)
	logger.LogInfo("CFG: Telegram BotId = ", BotId)
	logger.LogInfo("CFG: Telegram LogId = ", LogId)
	logger.LogInfo("CFG: Service host = ", Host)
	logger.LogInfo("CFG: Service base path = ", BasePath)
}
