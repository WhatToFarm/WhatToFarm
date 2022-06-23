package telegram

import (
	"ext-tg-bot/core"
	"ext-tg-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
)

const (
	basePath = "https://api.telegram.org/file/bot%s/%s"
)

func (bot *TgBot) handleFileUpload(message *tgbotapi.Message, userID int64) bool {
	fileName := message.Document.FileName
	fileConfig := tgbotapi.FileConfig{
		FileID: message.Document.FileID,
	}
	tgFile, err := bot.api.GetFile(fileConfig)
	if err != nil {
		bot.sendMessage(userID, wrong)
		utils.LogError("file data upload:", err)
		return false
	}

	path := "./assets/" + strconv.Itoa(int(userID)) + ".txt"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	endpoint := fmt.Sprintf(basePath, core.BotId, tgFile.FilePath)
	params := map[string]string{
		"Accept": message.Document.MimeType,
	}

	buf := make([]byte, 0)
	fileReader := tgbotapi.FileBytes{
		Name:  path,
		Bytes: buf,
	}

	res, err := bot.api.UploadFile(endpoint, params, fileName, fileReader)
	if err != nil {
		bot.sendMessage(userID, wrong)
		utils.LogError("file upload:", err, res.ErrorCode, res.Result)
		return false
	}

	fmt.Println(string(buf))
	//_, err = io.Copy(file, buf)
	//if err != nil {
	//	log.Println(err)
	//}
	//alerts := bot.getUserAlerts(userID)
	//sendAlertsList(bot, message.Chat.ID, alertsList, alerts...)
	return true
}
