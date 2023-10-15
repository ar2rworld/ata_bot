package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ar2rworld/ata_bot/app/myerror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BAN = "ban"

type TriggerCallbackQuery struct {
	CommandStruct
}

func NewTriggerCallbackQuery() *TriggerCallbackQuery {
	return &TriggerCallbackQuery{}
}

func (*TriggerCallbackQuery) GetName() string {
	return "TriggerCallbackQuery handler"
}

func (c *TriggerCallbackQuery) Exec(u *tgbotapi.Update) error {
	if u.CallbackQuery == nil {
		return nil
	}

	if ! c.Authorised(u.CallbackQuery.From.ID) {
		return nil
	}

	data := strings.Split(u.CallbackQuery.Data, "|,|")
	// required 3 tokens in data: <action>|,|<chatID>|,|<userData>|,|<newMemberMessageID>
	if len(data) != 4 {
		return nil
	}
	if data[0] != BAN {
		return nil
	}

	chatID, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		return myerror.NewError(fmt.Sprintf("error parsing int: %s", err))
	}
	userID, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return myerror.NewError(fmt.Sprintf("error parsing int: %s", err))
	}
	newMemberMessageID, err := strconv.Atoi(data[3])
	if err != nil {
		return myerror.NewError(fmt.Sprintf("error parsing int: %s", err))
	}

	user := &tgbotapi.User{ ID: userID }

	ataBot := *c.GetAtaBot()
	ataStorage := *c.GetStorage()

	err = ataBot.BanUser(chatID, user.ID, true)
	if err != nil {
		return err
	}

	banned, err := ataStorage.IsBanned(user)
	if err != nil {
		return err
	}
	if ! banned {
		err = ataStorage.AddToBanned(user)
		if err != nil {
			return err
		}	
	}
	
	// delete newMemberMessage
	err = ataBot.DeleteMessage(chatID, newMemberMessageID)

	return err
}
