package core

import (
	"github.com/spf13/viper"
	"ton-tg-bot/logger"
)

var (
	MongoURL      string
	MongoDatabase string
	LogLevel      string
	BotId         string
	LogId         int64
	Host          string
	Port          string
	BasePath      string
)

func Init() {
	MongoURL = viper.GetString("mongo.url")
	MongoDatabase = viper.GetString("mongo.database")
	LogLevel = viper.GetString("log.level")
	Host = viper.GetString("service.host")
	Port = viper.GetString("service.port")
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
	logger.LogInfo("CFG: LogLevel = ", LogLevel)
	logger.LogInfo("CFG: Telegram BotId = ", BotId)
	logger.LogInfo("CFG: Telegram LogId = ", LogId)
	logger.LogInfo("CFG: Service host = ", Host)
	logger.LogInfo("CFG: Service host = ", Port)
	logger.LogInfo("CFG: Service base path = ", BasePath)

	logger.LogLevel = LogLevel
}
