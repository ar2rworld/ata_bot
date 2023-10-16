package atabot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestGetUserProfilePic(t *testing.T) {
	ataBot := &AtaBot{}

	url, err := ataBot.GetUserProfilePic(&tgbotapi.User{
		UserName: "ne0ne0postaviat0kolonku",
	})
	if err != nil {
		t.Fatal(err)
	}
	if url == "" {
		t.Error("Empty profile pic returned")
	}
}