package telegram

import (
	"errors"
	"ext-tg-bot/core"
	"ext-tg-bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	basePath   = "https://api.telegram.org/file/bot%s/%s"
	defaultExt = ".txt"
)

var (
	InvalidFileFormat = errors.New("invalid file format")
)

func (bot *TgBot) handleFileUpload(message *tgbotapi.Message, userID int64) bool {
	err := bot.getFile(message)
	if errors.Is(err, InvalidFileFormat) {
		bot.sendMessage(userID, invalidFormat)
		return false
	} else if err != nil {
		bot.sendMessage(userID, invalidFormat)
		return false
	}

	return true
}

func (bot *TgBot) getFile(message *tgbotapi.Message) error {
	tgFile, err := bot.api.GetFile(tgbotapi.FileConfig{FileID: message.Document.FileID})
	if err != nil {
		utils.LogError("get file data:", err)
		return err
	}

	pathToFile, err := path(message)
	if err != nil {
		utils.LogError("get file data:", err)
		return err
	}

	file, err := os.Create(pathToFile)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if errFile := file.Close(); errFile != nil {
			utils.LogError("file close:", err)
		}
	}()

	url := fmt.Sprintf(basePath, core.BotId, tgFile.FilePath)
	return getFileFromURL(url, file)
}

func path(message *tgbotapi.Message) (string, error) {
	ext, err := validateFileExt(message.Document.FileName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("./assets/%d_%s%s", message.From.ID, time.Now().Format("06_01_02-15_04_05"), ext), err
}

func validateFileExt(fileName string) (string, error) {
	if fields := strings.Split(fileName, "."); len(fields) > 1 {
		if ext := fmt.Sprintf(".%s", fields[len(fields)-1]); ext == defaultExt {
			return ext, nil
		}
	}
	return "", InvalidFileFormat
}

func getFileFromURL(url string, file *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		utils.LogError("file download:", err)
		return err
	}
	defer func() {
		if errBody := resp.Body.Close(); errBody != nil {
			utils.LogError("body close:", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		utils.LogError("resp body copy:", err)
	}
	return err
}
