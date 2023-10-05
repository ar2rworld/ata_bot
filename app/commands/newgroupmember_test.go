package commands

import (
	"encoding/json"
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestCleanUserJoinsGroup(t *testing.T) {
	testUserID := int64(1)
	updateString := fmt.Sprintf(`{"update_id":382028976,
		"message":{"message_id":1508,"from":{"id":%d,"is_bot":false,"first_name":"Null","last_name":"User","language_code":"en"},
		"chat":{"id":-1001506079405,"title":"Public group for my bots","username":"my_bots_group","type":"supergroup"},
		"date":1696536356,
		"new_chat_participant":{"id":%d,"is_bot":false,"first_name":"Null","last_name":"User","language_code":"en"},
		"new_chat_member":{"id":%d,"is_bot":false,"first_name":"Null","last_name":"User","language_code":"en"},
		"new_chat_members":[{"id":%d,"is_bot":false,"first_name":"Null","last_name":"User","language_code":"en"}]}}`, testUserID, testUserID, testUserID, testUserID)
	var update *tgbotapi.Update
	err := json.Unmarshal([]byte(updateString), &update)
	if err != nil {
		t.Error(err)
	}

	var testStorage = &testStorage{}

	newGroupMember := NewNewGroupMember()
	newGroupMember.SetAtaBot(&ThisTestBot{})
	newGroupMember.SetStorage(testStorage)

	err = newGroupMember.Exec(update)
	if err != nil {
		t.Error(err)
	}

	if testStorage.targetUserID != testUserID {
		t.Errorf("targetUserID is incorrect: %d != %d", testStorage.targetUserID, testUserID)
	}

}

type testStorage struct {
	TestStorage
	targetUserID int64
}
func (t *testStorage) IsBanned(u *tgbotapi.User) (bool, error) {
	t.targetUserID = u.ID
	return false, nil
}