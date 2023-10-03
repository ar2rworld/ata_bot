package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewGroupMember struct {
	CommandStruct
}

func NewNewGroupMember () *NewGroupMember {
	return &NewGroupMember{}
}

func (n *NewGroupMember) GetName() string {
	return "NewGroupMember handler"
}

func (n *NewGroupMember) Exec(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if len(update.Message.NewChatMembers) > 0 {
		storage := *n.GetStorage()
		for _, chatMember := range update.Message.NewChatMembers {
			banned, err := storage.IsBanned(&chatMember)
			if err != nil {
				return err
			}
			if banned {

				err = n.ataBot.BanUser(update.Message.Chat.ID, chatMember.ID, true)

				if err != nil {
					return err
				}
			}
		}		
	}

	return nil
}
