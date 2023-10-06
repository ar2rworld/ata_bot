package commands

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewUpdate (in string) (*tgbotapi.Update, error) {
	var update = &tgbotapi.Update{}
	
	err := json.Unmarshal([]byte(in), update)
	
	return update, err
}
