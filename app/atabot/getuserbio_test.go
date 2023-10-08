package atabot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestGetUserBio(t *testing.T) {
	ataBot := &AtaBot{}
	user   := &tgbotapi.User{ UserName: "ne0ne0postaviat0kolonku" }

	bio, err := ataBot.GetUserBio(user)
	if err != nil {
		t.Error(err)
	}
	if bio != "dreams - no, goals - yes" {
		t.Errorf("Incorrect bio received OR admin changed bio: \"%s\"", bio)
	}
}