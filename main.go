package main

import (
	"context"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"ton-tg-bot/core"
	"ton-tg-bot/logger"
	"ton-tg-bot/mongo"
	"ton-tg-bot/telegram"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		logger.LogInfo("Using config file:", viper.ConfigFileUsed())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core.Init()
	initDb(ctx)

	logger.LogInfo("READY")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	tgBot := &telegram.TgBot{}
	tgBot.Init(core.BotId, core.LogId)
	go tgBot.Start()

	logger.LogInfo("Telegram bot started")

	for {
		select {
		case <-c:
			logger.LogInfo("Shutting down...")
			mongo.Close()
			os.Exit(0)
		}
	}
}

func initDb(ctx context.Context) {
	err := mongo.Init(ctx, core.MongoURL, core.MongoDatabase)
	if err != nil {
		logger.LogFatal("Mongo Init error:", err.Error())
	}
}
