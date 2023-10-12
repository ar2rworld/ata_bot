package commands

import (
	"fmt"
	"testing"

	"github.com/ar2rworld/ata_bot/app/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestTriggerIsBot(t *testing.T) {
	targetChatID := int64(1)
	targetUserID := int64(2)
	targetIsBot := true
	updateString := fmt.Sprintf(`{"update_id":382028976,
		"message":{"message_id":1508,"from":{"id":%d,"is_bot":%v,"first_name":"Null","last_name":"User","language_code":"en"},
		"chat":{"id":%d,"title":"Public group for my bots","username":"my_bots_group","type":"supergroup"},
		"date":1696536356,
		"new_chat_participant":{"id":%d,"is_bot":%v,"first_name":"Null","last_name":"User","language_code":"en"},
		"new_chat_member":{"id":%d,"is_bot":%v,"first_name":"Null","last_name":"User","language_code":"en"},
		"new_chat_members":[{"id":%d,"is_bot":%v,"first_name":"Null","last_name":"User","language_code":"en"}]}}`,
		targetUserID, targetIsBot, targetChatID, targetUserID, targetIsBot, targetUserID, targetIsBot, targetUserID, targetIsBot)

	update, err := NewUpdate(updateString)
	if err != nil {
		t.Error(err)
	}

	testBot := &triggerTestBot{}
	testStorage := &triggerTestStorage{}

	trigger := NewTrigger()
	trigger.SetAtaBot(testBot)
	trigger.SetStorage(testStorage)

	err = trigger.Exec(update)

	if err != nil {
		t.Error(err)
	}

	if testBot.bannedUserID != targetUserID {
		t.Errorf("bot: user id is incorrect: %d != %d", testBot.bannedUserID, targetUserID)
	}
	if testBot.bannedChatID != targetChatID {
		t.Errorf("bot: chat id is incorrect: %d != %d", testBot.bannedChatID, targetChatID)
	}

	if testStorage.bannedUserID != targetUserID {
		t.Errorf("storage: user id is incorrect: %d != %d", testStorage.bannedUserID, targetUserID)
	}
}

func TestTriggerUserDescription(t *testing.T) {
	// user joins the chat, bot has to check his bio and figure if the user has some interesting words in bio ban him

	testUserID := int64(1014210753)
	testChatID := int64(-1001506079405)
	testUserUserName := "a"
	updateString := fmt.Sprintf(`{"update_id":382028976,
		"message":{"message_id":1508,"from":{"id":%d,"is_bot":false,"username":"%s","first_name":"Null","last_name":"User","language_code":"en"},
		"chat":{"id":%d,"title":"Public group for my bots","username":"my_bots_group","type":"supergroup"},
		"date":1696536356,
		"new_chat_participant":{"id":%d,"is_bot":false,"first_name":"Null","last_name":"User","language_code":"en"},
		"new_chat_member":{"id":%d,"is_bot":false,"username":"%s","first_name":"Null","last_name":"User","language_code":"en"},
		"new_chat_members":[{"id":%d,"is_bot":false,"username":"%s","first_name":"Null","last_name":"User","language_code":"en"}]}}`,
		testUserID, testUserUserName, testChatID, testUserID, testUserID, testUserUserName, testUserID, testUserUserName)
	update, err := NewUpdate(updateString)
	if err != nil {
		t.Error(err)
	}
	
	var testStorage = &triggerDescriptionTestStorage{}
	var testBot     = &triggerDescriptionTestBot{}

	trigger := NewTrigger()
	trigger.SetAtaBot(testBot)
	trigger.SetStorage(testStorage)

	err = trigger.Exec(update)

	if err != nil {
		t.Fatal(err)
	}

	if testBot.bannedChatID != testChatID {
		t.Errorf("Incorrect bannedChatID: %d != %d", testBot.bannedChatID, testChatID)
	}
	if testBot.bannedUserID != testUserID {
		t.Errorf("Incorrect bannedChatID: %d != %d", testBot.bannedUserID, testUserID)
	}
	if testBot.bannedUserUserName != testUserUserName {
		t.Errorf("Incorrect bannedUserUserName: %s != %s", testBot.bannedUserUserName, testUserUserName)
	}
	if testStorage.bannedUserID != testUserID {
		t.Errorf("Storage: Incorrect bannedUserID: %d != %d", testStorage.bannedUserID, testUserID)
	}

}

type triggerTestBot struct {
	TestBot
	bannedChatID int64
	bannedUserID int64
}
func (b *triggerTestBot) BanUser(chatID, userID int64, revokeMessages bool) error {
	b.bannedChatID = chatID
	b.bannedUserID = userID
	return nil
}

type triggerTestStorage struct {
	TestStorage
	bannedUserID int64
}
func (s *triggerTestStorage) AddToBanned(u *tgbotapi.User) error {
	s.bannedUserID = u.ID
	return nil
}

type triggerDescriptionTestStorage struct {
	TestStorage
	bannedUserID int64
}
func (s *triggerDescriptionTestStorage) GetTriggerWords() (*[]storage.TriggerWord, error) {
	return &[]storage.TriggerWord{
		{
			Text: "asdf",
			Severity: storage.Severity200,
		},
	}, nil
}
func (s *triggerDescriptionTestStorage) AddToBanned(u *tgbotapi.User) error {
	s.bannedUserID = u.ID
	return nil
}

type triggerDescriptionTestBot struct {
	TestBot
	bannedChatID int64
	bannedUserID int64
	bannedUserUserName string
}
func (b *triggerDescriptionTestBot) GetUserBio (u *tgbotapi.User) (string, error) {
	b.bannedUserUserName = u.UserName
	return "something asdf smt", nil
}
func (b *triggerDescriptionTestBot) BanUser(chatID, userID int64, revokeMessages bool) error {
	b.bannedChatID = chatID
	b.bannedUserID = userID
	return nil
}