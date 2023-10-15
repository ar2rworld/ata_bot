package commands

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestTriggerCallbackQuery(t *testing.T) {
	t.Run("Check Authorised command", func(t *testing.T) {
		adminID := int64(1)
		unauthorisedID := int64(2)
		update := &tgbotapi.Update{
			CallbackQuery: &tgbotapi.CallbackQuery{
				From: &tgbotapi.User{ ID: unauthorisedID },
				Data: BAN + "|,|123|,|321",
			},
		}

		tcq := NewTriggerCallbackQuery()

		ataBot := &triggerCallbackQueryTestAuthorizedBot{}
		ataStorage := &TestStorage{}

		tcq.SetAtaBot(ataBot)
		tcq.SetStorage(ataStorage)

		err := tcq.Exec(update)
		if err != nil {
			t.Fatal(err)
		}

		if ! tcq.Authorised(adminID) {
			t.Error("AdminID should be set to the command and thus authorized")
		}

		if tcq.Authorised(unauthorisedID) {
			t.Error("UnauthorisedID should not be authorized")
		}

		if ataBot.accessedBanUser {
			t.Errorf("Accessed by unauthorised user with id: %d", update.CallbackQuery.From.ID)
		}
	})

	t.Run("Successful triggerCallbackQuery", func(t *testing.T) {
		userID := int64(1014210753)
		chatID := int64(1014210753)
		
		action       := BAN
		publicChatID := int64(-1001506079405)
		targetUserID := int64(1265820975)
		

		updateString := fmt.Sprintf(`{"update_id":382029020,
		"callback_query":{
			"id":"4356002017066781719",
			"from":{
				"id":%d,"is_bot":false,"first_name":"Nemo","last_name":"Cap","username":"ne0ne0postaviat0kolonku","language_code":"en"
			},
			"message":{
				"message_id":3666,
				"from":{"id":1832953211,"is_bot":true,"first_name":"NemoBotPhone","username":"nemobotphone_bot"},
				"chat":{
					"id":%d,"first_name":"Nemo","last_name":"Cap","username":"ne0ne0postaviat0kolonku","type":"private"
				},
				"date":1697221643,"text":"suspicious user(1265820975) bio: \"asdf\" in chat(-1001506079405)",
				"reply_markup":{
					"inline_keyboard":[[{"text":"ban user","callback_data":"ban,-1001506079405,1265820975"},{"text":"profile","url":"https://t.me/whybotheraboutusername"}]]
				}
				},
				"chat_instance":"-8603891193044155368",
				"data":"%s|,|%d|,|%d"
			}
		}`, userID, chatID, action, publicChatID, targetUserID)
		update, err := NewUpdate(updateString)
		if err != nil {
			t.Fatal(err)
		}

		testBot := &triggerCallbackQueryTestBot{}
		testStorage := &triggerCallbackQueryTestStorage{}

		triggerCallbackQuery := NewTriggerCallbackQuery()
		triggerCallbackQuery.SetAtaBot(testBot)
		triggerCallbackQuery.SetStorage(testStorage)

		err = triggerCallbackQuery.Exec(update)

		if err != nil {
			t.Error(err)
		}
		if testBot.bannedChatID != publicChatID {
			t.Errorf("Incorrect bannedChatID: %d != %d", testBot.bannedChatID, publicChatID)
		}
		if testBot.bannedUserID != targetUserID {
			t.Errorf("Incorrect bannedUserID: %d != %d", testBot.bannedUserID, targetUserID)
		}
		if ! testBot.revokeMessage {
			t.Error("RevokeMessages should be true")
		}
		if testStorage.isBannedUserID != targetUserID {
			t.Errorf("Storage: isBannedUserID incorrect: %d != %d", testStorage.isBannedUserID, targetUserID)
		}
		if testStorage.bannedUserID != targetUserID {
			t.Errorf("Storage: bannedUserID incorrect: %d != %d", testStorage.bannedUserID, targetUserID)
		}

	})
}

type triggerCallbackQueryTestAuthorizedBot struct {
	TestBot
	accessedBanUser bool
}
func (*triggerCallbackQueryTestAuthorizedBot) GetAdminID() int64 {
	return 1
}
func (b *triggerCallbackQueryTestAuthorizedBot) BanUser(_, _ int64, _ bool) error {
	b.accessedBanUser = true
	return nil
}

type triggerCallbackQueryTestBot struct {
	TestBot
	bannedChatID int64
	bannedUserID int64
	revokeMessage bool
}
func (b *triggerCallbackQueryTestBot) BanUser(chatID, userID int64, revokeMessage bool) error {
	b.bannedChatID = chatID
	b.bannedUserID = userID
	b.revokeMessage = revokeMessage
	return nil
}
func (*triggerCallbackQueryTestBot) GetAdminID() int64 {
	return 1014210753
}

type triggerCallbackQueryTestStorage struct {
	TestStorage
	isBannedUserID int64
	bannedUserID int64
}

func (s *triggerCallbackQueryTestStorage) IsBanned(u *tgbotapi.User) (bool, error) {
	s.isBannedUserID = u.ID
	return false, nil
}

func (s *triggerCallbackQueryTestStorage) AddToBanned(u *tgbotapi.User) error {
	s.bannedUserID = u.ID
	return nil
}
