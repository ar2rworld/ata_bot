package commands

import (
	"fmt"
	"strings"

	"github.com/ar2rworld/ata_bot/app/myerror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const apistats = "/apistats"

type APIStats struct {
	CommandStruct
}

func NewAPIStats () *APIStats {
	return &APIStats{}
}

func (h *APIStats) GetName() string {
	return apistats
}

func (c *APIStats) Exec(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if update.Message.Text != apistats {
		return nil
	}
	if ! c.Authorised(update.Message.From.ID) {
		return myerror.NewError("is not admin")
	}

	var out = []string{ "Header:\n" }
	ataBot := *c.GetAtaBot()
	header := ataBot.GetAPIResponseHeader()
	for _, h := range header {
		if len(h) != 2 {
			return myerror.NewError("Invalid header: " + strings.Join(h, ","))
		}
		out = append(out, fmt.Sprintf("%s: %s", h[0], h[1]))
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(out, "\n"))
	_, err := ataBot.Send(message)

	return err
}
