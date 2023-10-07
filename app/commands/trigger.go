package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Trigger struct {
	CommandStruct
}

func NewTrigger() *Trigger {
	return &Trigger{}
}

func (n *Trigger) GetName() string {
	return "Trigger"
}
func (n *Trigger) GetHelp() string {
	return "when new member joins the chat, bot will check if user is bot or has some trigger words when ban him"
}

func (t *Trigger) Exec(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if len(update.Message.NewChatMembers) > 0 {
		ataBot  := *t.GetAtaBot()
		storage := *t.GetStorage()

		for _, newMember := range update.Message.NewChatMembers {
			
			// if IsBot
			if newMember.IsBot {
				err := ataBot.BanUser(update.Message.Chat.ID, newMember.ID, true)
				if err != nil {
					return err
				}
				err = storage.AddToBanned(&newMember)
				if err != nil {
					return nil
				}
			}

			// TODO: if newMember has some interesting words in bio
		}

	}
	return nil
}
