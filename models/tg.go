package models

const (
	StateInitial = iota
	StateMessage
)

type TgUser struct {
	TgId  int64
	State int
}
