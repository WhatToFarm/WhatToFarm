package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"ton-tg-bot/core"
	"ton-tg-bot/logger"
)

const (
	basePath   = "https://api.telegram.org/file/bot%s/%s"
	defaultExt = ".tar.gz"
)

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
	ErrInternal          = errors.New("internal error")
	ErrToManyAttempts    = errors.New("many attempts")
)

func (bot *TgBot) handleFileUpload(message *tgbotapi.Message, userID int64) {
	if err := bot.checkAttempt(userID); err != nil {
		return
	}

	err := bot.getFile(message)
	if errors.Is(err, ErrInvalidFileFormat) {
		bot.sendMessage(userID, invalidFormat)
		return
	} else if err != nil {
		bot.sendMessage(userID, wrong)
		return
	}
}

func (bot *TgBot) checkAttempt(userID int64) error {
	user, ok := bot.users[userID]
	if !ok {
		bot.sendMessage(userID, wrong)
		logger.LogError("user not added to bot cache", userID)
		return ErrInternal
	}

	timeLimit := time.Now().UTC().Add(-1 * time.Hour)
	if user.Attempts == 5 && user.TS.After(timeLimit) {
		next := timeLimit.Unix() - user.TS.Unix()
		bot.sendMessage(userID, func() string {
			return toManyAttempts(next)
		})
		return ErrToManyAttempts
	}

	if user.TS.After(timeLimit) {
		user.Attempts++
	} else {
		user.Attempts = 1
		user.TS = time.Now().UTC()
	}

	return nil
}

func (bot *TgBot) getFile(message *tgbotapi.Message) error {
	tgFile, err := bot.api.GetFile(tgbotapi.FileConfig{FileID: message.Document.FileID})
	if err != nil {
		logger.LogError("get file data:", err)
		return err
	}

	pathToFile, err := path(message)
	if err != nil {
		logger.LogError("get file data:", err)
		return err
	}

	file, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.LogError("file create:", err)
		return err
	}
	defer func() {
		if errFile := file.Close(); errFile != nil {
			logger.LogError("file close:", err)
		}
	}()

	url := fmt.Sprintf(basePath, core.BotId, tgFile.FilePath)
	return getFileFromURL(url, file)
}

func path(message *tgbotapi.Message) (string, error) {
	if !strings.HasSuffix(message.Document.FileName, defaultExt) {
		return "", ErrInvalidFileFormat
	}
	return fmt.Sprintf("./assets/%d_%s%s", message.From.ID, time.Now().Format("0102150405"), defaultExt), nil
}

func getFileFromURL(url string, file *os.File) error {
	resp, err := http.Get(url)
	if err != nil {
		logger.LogError("file download:", err)
		return err
	}
	defer func() {
		if errBody := resp.Body.Close(); errBody != nil {
			logger.LogError("body close:", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		logger.LogError("resp body copy:", err)
	}
	return err
}
