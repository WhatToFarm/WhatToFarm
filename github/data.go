package github

import (
	"errors"
	"time"
)

const (
	pathUser = "https://api.github.com/users/%s"
	pathRepo = "https://api.github.com/repos/%s/TCS2"

	month = time.Hour * 24 * 30
)

var (
	ErrNotFound = errors.New("GitHub user not found")
	ErrNewUser  = errors.New("GitHub user registered less than a month")
)
