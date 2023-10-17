package commands

import (
	"fmt"
	"testing"

	"github.com/ar2rworld/ata_bot/app/api"
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
	
	var testStorage = &triggerDescriptionTestStorage{}
	var testBot     = &triggerDescriptionTestBot{}

	trigger := NewTrigger()
	trigger.SetAtaBot(testBot)
	trigger.SetStorage(testStorage)

	t.Run("Ban user with interesting bio", func(t *testing.T) {
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
			t.Fatal(err)
		}
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
		if testStorage.report.ChatID != testChatID ||
		testStorage.report.UserID != testUserID ||
		testStorage.report.Severity != storage.Severity200 ||
		testStorage.report.Action != storage.ActionBanned {
			t.Errorf("Something is wrong with Report: %v", testStorage.report)
		}
		if testBot.sendIsCalled {
			t.Error("SendToAdmin should not be called")
		}
	})
	t.Run("Notify admin about user with sus bio", func(t *testing.T) {
		testUserID := int64(1)
		testChatID := int64(-1)
		testUserUserName := "b"
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

		err = trigger.Exec(update)

		if testBot.bannedChatID == testChatID {
			t.Errorf("Incorrect bannedChatID: %d != %d", testBot.bannedChatID, testChatID)
		}
		if testBot.bannedUserID == testUserID {
			t.Errorf("Incorrect bannedChatID: %d != %d", testBot.bannedUserID, testUserID)
		}
		if testBot.bannedUserUserName == testUserUserName {
			t.Errorf("Incorrect bannedUserUserName: %s != %s", testBot.bannedUserUserName, testUserUserName)
		}
		if testStorage.bannedUserID == testUserID {
			t.Errorf("Storage: Incorrect bannedUserID: %d != %d", testStorage.bannedUserID, testUserID)
		}
		if testStorage.report.ChatID != testChatID ||
		testStorage.report.UserID != testUserID ||
		testStorage.report.Severity != storage.Severity100 ||
		testStorage.report.Action != "notified" {
			t.Errorf("Something is wrong with Report: %v", testStorage.report)
		}
		if ! testBot.sendIsCalled {
			t.Error("Send should be called")
		}
		if testBot.AnalyzeUserPicIsCalled {
			t.Error("AnalyzeUserPicIsCalled should not be called")
		}
	})
}

func TestTriggerAnalyzeUserPic(t *testing.T) {
	testBot := &triggerAnalyzeUserPicTestBot{}
	testStorage := &triggerAnalyzeUserPicTestStorage{}

	trigger := NewTrigger()
	trigger.SetAtaBot(testBot)
	trigger.SetStorage(testStorage)

	t.Run("Unsafe user pic", func(t *testing.T) {
		testUserID := int64(1014210753)
		testChatID := int64(-1001506079405)
		testUserUserName := "c"
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
			t.Fatal(err)
		}

		if err = trigger.Exec(update); err != nil {
			t.Fatal(err)
		}

		if testBot.analyzedNewMemberID != testUserID {
			t.Errorf("Incorrect analyzedNewMemberID: %d != %d", testBot.analyzedNewMemberID, testUserID)
		}
		if testBot.analyzedNewMemberUserName != testUserUserName {
			t.Errorf(`Incorrect analyzedNewMemberUserName: "%s" != "%s"`, testBot.analyzedNewMemberUserName, testUserUserName)
		}
		if testBot.sentMessage == nil {
			t.Error("bot should have sent a message")
		}
		if testStorage.report.Comment != "unsafe profile pic" {
			t.Error("report comment should be \"unsafe profile pic\":" , testStorage.report.Comment)
		}
	})

	t.Run("User pic is not unsafe but potentially unsafe", func(t *testing.T) {
		testUserID := int64(123)
		testChatID := int64(-456)
		testUserUserName := "d"
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
			t.Fatal(err)
		}

		if err = trigger.Exec(update); err != nil {
			t.Fatal(err)
		}

		if testBot.analyzedNewMemberID != testUserID {
			t.Errorf("Incorrect analyzedNewMemberID: %d != %d", testBot.analyzedNewMemberID, testUserID)
		}
		if testBot.analyzedNewMemberUserName != testUserUserName {
			t.Errorf(`Incorrect analyzedNewMemberUserName: "%s" != "%s"`, testBot.analyzedNewMemberUserName, testUserUserName)
		}
		if testBot.sentMessage == nil {
			t.Error("bot should have sent a message")
		}
		if testStorage.report.Comment != "unsafe profile pic" {
			t.Error("report comment should be \"unsafe profile pic\":" , testStorage.report.Comment)
		}
	})
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
	report storage.ReportStruct
}
func (s *triggerDescriptionTestStorage) GetTriggerWords() (*[]storage.TriggerWord, error) {
	return &[]storage.TriggerWord{
		{
			Text: "asdf",
			Severity: storage.Severity200,
		},
		{
			Text: "fdsa",
			Severity: storage.Severity100,
		},
	}, nil
}
func (s *triggerDescriptionTestStorage) AddToBanned(u *tgbotapi.User) error {
	s.bannedUserID = u.ID
	return nil
}
func (s *triggerDescriptionTestStorage) Report(chatID, userID int64, severity int, action, comment string) error {
	s.report.ChatID = chatID
	s.report.UserID = userID
	s.report.Severity = severity
	s.report.Action = action
	s.report.Comment = comment
	return nil
}
type triggerDescriptionTestBot struct {
	TestBot
	bannedChatID int64
	bannedUserID int64
	bannedUserUserName string
	sendIsCalled bool
	AnalyzeUserPicIsCalled bool
}
func (b *triggerDescriptionTestBot) GetUserBio (u *tgbotapi.User) (string, error) {
	if u.UserName == "a" {
		b.bannedUserUserName = u.UserName
		return "something asdf smt", nil
	}
	return "fdsa", nil
}
func (b *triggerDescriptionTestBot) BanUser(chatID, userID int64, revokeMessages bool) error {
	if userID == int64(1014210753) && chatID == int64(-1001506079405) {
		b.bannedChatID = chatID
		b.bannedUserID = userID
	}
	return nil
}
func (t *triggerDescriptionTestBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	t.sendIsCalled = true
	return tgbotapi.Message{}, nil
}
func (t *triggerDescriptionTestBot) AnalyzeUserPic(u *tgbotapi.User) (api.APIResponse, error) {
	t.AnalyzeUserPicIsCalled = true
	return api.APIResponse{Unsafe: true}, nil
}

type triggerAnalyzeUserPicTestBot struct {
	TestBot
	sentMessage tgbotapi.Chattable
	analyzedNewMemberID int64
	analyzedNewMemberUserName string
}

func (t *triggerAnalyzeUserPicTestBot) AnalyzeUserPic(u *tgbotapi.User) (api.APIResponse, error) {
	if u.UserName == "c" {
		t.analyzedNewMemberID = u.ID
		t.analyzedNewMemberUserName = u.UserName
		return api.APIResponse{Unsafe: true}, nil	
	} else if u.UserName == "d" {
		t.analyzedNewMemberID = u.ID
		t.analyzedNewMemberUserName = u.UserName
		return api.APIResponse{
			Unsafe: false,
			Objects: []api.Objects{
				{
					Label: api.EXPOSED_ANUS,
					Score: float32(0.5),
				},
			},
		}, nil	
	}
	return api.APIResponse{}, nil	
}

func (t *triggerAnalyzeUserPicTestBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	t.sentMessage = c
	return tgbotapi.Message{}, nil
}

type triggerAnalyzeUserPicTestStorage struct {
	TestStorage
	report storage.ReportStruct
}

func (s *triggerAnalyzeUserPicTestStorage) Report(chatID int64, userID int64, severity int, action string, comment string) error {
	s.report = *storage.NewReport(chatID, userID, severity, action, comment)
	return nil
}
