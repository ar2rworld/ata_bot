package commands

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestWatchReply(t *testing.T) {
	bannedChatID := int64(-1001845507313)
	bannedBotID := int64(260566685)
	bannedUserName := "chelpbot"
	updateString := fmt.Sprintf(`{"update_id":127504112,
		"message":{"message_id":2650,"from":{"id":6618529494,"is_bot":false,"first_name":"Deborah"},
		"chat":{"id":%d,"title":"Name of the group in St. John's, NL","username":"nameofthegroupinstjohnsnl","type":"supergroup"},
		"date":1698294833,
		"video":{"duration":6,"width":270,"height":320,"mime_type":"video/mp4",
			"thumbnail":{"file_id":"tgfileid-tgfileid","file_unique_id":"tgfileid","file_size":14207,"width":270,"height":320},
			"thumb":{"file_id":"tgfileid-tgfileid","file_unique_id":"tgfileid","file_size":14207,"width":270,"height":320},
			"file_id":"tgfileid-tgfileid","file_unique_id":"tgfileid","file_size":137701},
			"reply_markup":{"inline_keyboard":[
				[{"text":"\ud83d\udd25 some string \ud83d\udd25","url":"someurl"}],
				[{"text":"\ud83d\udc95 some string \ud83d\udc95","url":"someurl"}],
				[{"text":"\ud83d\udca5 some string \ud83d\udca5","url":"someurl"}],
				[{"text":"\ud83d\udd25 some string \ud83d\udd25","url":"someurl"}],
				[{"text":"\ud83d\udc9e some string \ud83d\udc9e","url":"someurl"}],
				[{"text":"\ud83d\udc9d some string \ud83d\udc9d","url":"someurl"}]
			]},
		"via_bot":{"id":%d,"is_bot":true,"first_name":"Channel Help","username":"%s"}}}`,
		bannedChatID, bannedBotID, bannedUserName,
	)

	u, err := NewUpdate(updateString)
	if err != nil {
		t.Fatal(err)
	}

	testBot := &watchReplyTestBot{}
	testStorage := &watchReplyTestStorage{}

	c := NewWatchReply()
	c.SetAtaBot(testBot)
	c.SetStorage(testStorage)

	err = c.Exec(u)
	if err != nil {
		t.Fatal(err)
	}

	if testBot.bannedChatID != bannedChatID {
		t.Errorf("Incorrect bannedChatID: %d != %d", testBot.bannedChatID, bannedChatID)
	}
	if testBot.bannedUserID != bannedBotID {
		t.Errorf("Incorrect bannedUserID: %d != %d", testBot.bannedUserID, bannedBotID)
	}
	if testStorage.isBannedUserID != bannedBotID {
		t.Errorf("Incorrect isBannedUserID: %d != %d", testStorage.isBannedUserID, bannedBotID)
	}
	if testStorage.isBannedUserName != bannedUserName {
		t.Errorf("Incorrect isBannedUserName: %s != %s", testStorage.isBannedUserName, bannedUserName)
	}
}

type watchReplyTestBot struct {
	TestBot
	bannedChatID int64
	bannedUserID int64
}
func (b *watchReplyTestBot) BanUser(chatID, userID int64, revokeMessages bool) error {
	b.bannedChatID = chatID
	b.bannedUserID = userID
	return nil
}

type watchReplyTestStorage struct {
	TestStorage
	isBannedUserID int64
	isBannedUserName string
}
func (s *watchReplyTestStorage) IsBanned(u *tgbotapi.User) (bool, error) {
	s.isBannedUserID = u.ID
	s.isBannedUserName = u.UserName
	return true, nil
}