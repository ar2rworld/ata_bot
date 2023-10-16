package atabot

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/ar2rworld/ata_bot/app/myerror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *AtaBot) GetUserProfilePic(u *tgbotapi.User) (string, error) {
	if u.UserName == "" {
		return "", myerror.NewError("user missing username: " + u.String())
	}

	// request telegram page and parse possible bio
	res, err := http.Get(fmt.Sprintf("https://t.me/%s", u.UserName))
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	page := string(data)

	// page contains div with class "tgme_page_photo_image"
	re := regexp.MustCompile(`<img\s*class="tgme_page_photo_image\s*"\s*src="(.*)"`)
	matches := re.FindStringSubmatch(page)
	if len(matches) < 2 {
		return "", myerror.NewError(fmt.Sprintf("Could not find profile pic, %d, %s", len(matches), matches))
	}
	
	return matches[1], nil
}