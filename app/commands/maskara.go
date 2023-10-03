package commands

import (
	"strings"

	"github.com/ar2rworld/ata_bot/app/myerror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maskara = "/maskara"

type Maskara struct {
	CommandStruct
}

func NewMaskara() *Maskara {
	return &Maskara{}
}

func (m *Maskara) GetName() string {
	return maskara
}

func (m *Maskara) Exec(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if update.Message.ReplyToMessage == nil {
		return nil
	}

	hasNewChatMemberInReply := len(update.Message.ReplyToMessage.NewChatMembers) > 0

	messageText := update.Message.Text
	if _, found := strings.CutPrefix(messageText, maskara); found {
		if ! m.Authorised(update.Message.From.ID) {
			return myerror.NewError("is not admin")
		}

		storage := *m.GetStorage()

		var fakeUser *tgbotapi.User
		if hasNewChatMemberInReply {
			fakeUser = &update.Message.ReplyToMessage.NewChatMembers[0]
		} else {
			fakeUser = update.Message.ReplyToMessage.From
		}
		isBanned, err := storage.IsBanned(fakeUser)
		if err != nil && err.Error() != myerror.NoDocuments {
			return err
		}

		if ! isBanned {
			err = storage.AddToBanned(fakeUser)
			if err != nil {
				return err
			}
		}

		ataBot := *m.GetAtaBot()
		chatID := update.Message.Chat.ID

		err = ataBot.BanUser(chatID, fakeUser.ID, true)
		if err != nil && err.Error() == myerror.BadRequestNotEnoughRights {
			askForPermissions := tgbotapi.NewMessage(chatID, "need \"Ban Users\" and \"Delete Messages\" permission")
			ataBot.Send(askForPermissions)
		} else if err != nil && err.Error() == myerror.UnmarshallBoolErrorMessage {
			// do nothing
		} else {
			return err
		}

		clearMaskara := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
		_, err = ataBot.Send(clearMaskara)
		if err != nil && err.Error() != myerror.UnmarshallBoolErrorMessage {
			return err
		}

		clearFake := tgbotapi.NewDeleteMessage(chatID, update.Message.ReplyToMessage.MessageID)
		_, err = ataBot.Send(clearFake)
		if err != nil && err.Error() != myerror.UnmarshallBoolErrorMessage {
			return err
		}
	}
	return nil
}
