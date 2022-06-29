package models

import "time"

type ErrMessage string

const (
	NotFound ErrMessage = "Not Found"

	FieldAccount  = "gitAccount"
	FieldTS       = "ts"
	FieldAttempts = "attempts"
)

type Data struct {
	ID        int        `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Message   ErrMessage `json:"message"`
}

type TgUser struct {
	TgId       int64     `bson:"tgId"`
	State      int       `bson:"-"`
	TS         time.Time `bson:"ts"`
	Attempts   int       `bson:"attempts"`
	GitAccount string    `bson:"gitAccount"`
}
