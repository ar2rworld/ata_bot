package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewGroupUpdate struct {
	CommandStruct
}

func NewNewGroupUpdate () *NewGroupUpdate {
	return &NewGroupUpdate{}
}

func (n *NewGroupUpdate) GetName() string {
	return "NewGroupUpdate handler"
}

func (n *NewGroupUpdate) Exec(update *tgbotapi.Update) error {
	var err error
	if update.MyChatMember != nil {

		storage := *n.GetStorage()
		status := update.MyChatMember.NewChatMember.Status
		switch status {
			case "member":
				err = storage.AddNewChat(&update.MyChatMember.Chat, status)
				break
			case "restricted":
				err = storage.RemoveChat(&update.MyChatMember.Chat, status)
				break
			case "left":
				err = storage.RemoveChat(&update.MyChatMember.Chat, status)
				break
			case "kicked":
				err = storage.RemoveChat(&update.MyChatMember.Chat, status)
				break
		}

	}

	return err
}
