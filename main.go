package main

import (
	"context"
	"ext-tg-bot/core"
	"ext-tg-bot/mongo"
	"ext-tg-bot/telegram"
	"ext-tg-bot/utils"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		utils.LogInfo("Using config file:", viper.ConfigFileUsed())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core.Init()
	initDb(ctx)

	utils.LogInfo("READY")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	tgBot := &telegram.TgBot{}
	go tgBot.Init(core.BotId, core.LogId)

	utils.LogInfo("Telegram bot started")

	for {
		select {
		case <-c:
			utils.LogInfo("Shutting down...")
			mongo.Close()
			os.Exit(0)
		}
	}
}

func initDb(ctx context.Context) {
	err := mongo.Init(ctx, core.MongoURL, core.MongoDatabase)
	if err != nil {
		utils.LogError("Mongo Init error:", err.Error())
	}
}
