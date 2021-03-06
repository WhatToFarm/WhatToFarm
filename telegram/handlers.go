package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"ton-tg-bot/core"
	"ton-tg-bot/external"
	"ton-tg-bot/logger"
	"ton-tg-bot/models"
	"ton-tg-bot/mongo"
)

const (
	tgApiGetFile = "https://api.telegram.org/file/bot%s/%s"
	defaultExt   = ".tar.gz"

	maxFileSie = 1_000_000
)

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
	ErrManyAttempts      = errors.New("many attempts")
	ErrBigSize           = errors.New("big file size")
)

func (bot *TgBot) handleFileUpload(message *tgbotapi.Message, user *models.TgUser) {
	if !bot.checkDeadline(user) {
		return
	}

	if err := bot.checkAttempt(user); err != nil {
		return
	}

	resp, err := bot.handleFile(message)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidFileFormat):
			logger.LogWarn(err)
			bot.sendMessage(user.TgId, invalidFormat)
		case errors.Is(err, ErrBigSize):
			logger.LogWarn(err)
			bot.sendMessage(user.TgId, bigFileSize)
		default:
			logger.LogError(err)
			bot.sendMessage(user.TgId, wrong)
		}
		return
	}

	if err = mongo.UpdateUser(user); err != nil {
		logger.LogWarn("Update user attempts:", err)
	}

	bot.logText(message, "Response from service:\n"+resp)
	bot.sendMessage(user.TgId, func() string { return resp })
}

// checkDeadline - controls deadline.
func (bot *TgBot) checkDeadline(user *models.TgUser) bool {
	if time.Now().UTC().After(core.Deadline) {
		bot.sendMessage(user.TgId, expireDeadline)
		return false
	}
	return true
}

// checkAttempt - controls count of attempts.
// Should be 5 attempts per one hour.
func (bot *TgBot) checkAttempt(user *models.TgUser) error {
	timeLimit := time.Now().UTC().Add(-1 * time.Hour)
	if user.Attempts == 5 && user.TS.After(timeLimit) {
		next := (user.TS.Unix() - timeLimit.Unix()) / 60
		bot.sendMessage(user.TgId, func() string { return manyAttempts(next) })
		return ErrManyAttempts
	}

	if user.TS.After(timeLimit) {
		user.Attempts++
	} else {
		user.Attempts = 1
		user.TS = time.Now().UTC()
	}

	return nil
}

func (bot *TgBot) handleFile(message *tgbotapi.Message) (string, error) {
	tgFile, err := bot.api.GetFile(tgbotapi.FileConfig{FileID: message.Document.FileID})
	if err != nil {
		return "", fmt.Errorf("get file data: %w", err)
	}
	if tgFile.FileSize > maxFileSie {
		return "", ErrBigSize
	}

	pathToFile, err := path(message)
	if err != nil {
		return "", fmt.Errorf("file path: %w", err)
	}

	file, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("file create: %w", err)
	}
	defer func() {
		if errFile := file.Close(); errFile != nil {
			logger.LogError("file close:", err)
		}
	}()

	url := fmt.Sprintf(tgApiGetFile, core.BotId, tgFile.FilePath)
	if err = getFileFromURL(url, file); err != nil {
		return "", fmt.Errorf("get file from telegramm: %w", err)
	}

	message.Text = filepath.Base(pathToFile)
	bot.logForward(message)
	bot.sendMessage(message.Chat.ID, waitResult)

	buf := bytes.NewBufferString("")
	if _, err = io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("file copy: %w", err)
	}
	ans, err := external.ExtService(filepath.Base(pathToFile), buf)
	if err != nil {
		return "", fmt.Errorf("external service: %w", err)
	}
	return ans, nil
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
		return fmt.Errorf("file download: %w", err)
	}
	defer func() {
		if errBody := resp.Body.Close(); errBody != nil {
			logger.LogError("body close:", err)
		}
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("resp body copy: %w", err)
	}

	return nil
}
