package commands

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const help = "/help"

type Help struct {
	CommandStruct
}

func NewHelp () *Help {
	return &Help{}
}

func (h *Help) GetName() string {
	return help
}

func (n *Help) Exec(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if update.Message.Text != help {
		return nil
	}

	var out = []string{ "Commands:" }
	ataBot := *n.GetAtaBot()
	for _, c := range ataBot.GetCommands() {
		out = append(out, fmt.Sprintf(" - %s", c.GetName()))
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(out, "\n"))
	_, err := ataBot.Send(message)

	return err
}
