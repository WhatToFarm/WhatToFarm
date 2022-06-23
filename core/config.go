package core

import (
	"ext-tg-bot/utils"
	"github.com/spf13/viper"
)

var (
	MongoURL      string
	MongoDatabase string
	LogLevel      string
	BotId         string
	LogId         int64
)

func Init() {
	MongoURL = viper.GetString("mongo.url")
	MongoDatabase = viper.GetString("mongo.database")
	LogLevel = viper.GetString("log.level")

	BotId = viper.GetString("telegram.bot_id")
	if BotId == "" {
		utils.LogFatal("Telegram bot ID not provided")
	}

	LogId = viper.GetInt64("telegram.log_id")
	if LogId == 0 {
		utils.LogFatal("Telegram log ID not provided")
	}

	utils.LogInfo("CFG: MongoURL = ", MongoURL)
	utils.LogInfo("CFG: MongoDatabase = ", MongoDatabase)
	utils.LogInfo("CFG: LogLevel = ", LogLevel)
	utils.LogInfo("CFG: Telegram BotId = ", BotId)
	utils.LogInfo("CFG: Telegram LogId = ", LogId)

	utils.LogLevel = LogLevel
}
