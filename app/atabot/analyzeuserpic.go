package atabot

import (
	"strings"
	"net/http"
	"io"
	"os"
	"fmt"

	"github.com/ar2rworld/ata_bot/app/api"
	"github.com/ar2rworld/ata_bot/app/myerror"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *AtaBot) AnalyzeUserPic(u *tgbotapi.User) (api.APIResponse, error) {
	var voidAPIResponse api.APIResponse
	userPic, err := b.GetUserProfilePic(u)
	if err != nil {
		return voidAPIResponse, err
	}

	url := "https://nsfw-images-detection-and-classification.p.rapidapi.com/adult-content"

	payload := strings.NewReader(fmt.Sprintf("{\n\"url\": \"%s\"\n}", userPic))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return voidAPIResponse, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("X_RAPIDAPI_KEY"))
	req.Header.Add("X-RapidAPI-Host", "nsfw-images-detection-and-classification.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return voidAPIResponse, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return voidAPIResponse, err
	}

	if res.StatusCode != 200 {
		return voidAPIResponse, myerror.NewError(string(body))
	}

	b.SetAPIResponseHeader(&res.Header)

	apires, err := api.ParseAPI(body)
	if err != nil {
		return voidAPIResponse, err
	}
	
	return apires, nil
}
