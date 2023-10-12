package commands

import (
	"github.com/ar2rworld/ata_bot/app/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TestStorage struct {}

func (t *TestStorage) IsBanned(*tgbotapi.User) (bool, error) {
	return false, nil
}
func (t *TestStorage) AddToBanned(u *tgbotapi.User) error {
	return nil
}
func (t *TestStorage) AddNewChat(c *tgbotapi.Chat, status string) error {
	return nil
}
func (t *TestStorage) RemoveChat(c *tgbotapi.Chat, status string) error {
	return nil
}
func (t *TestStorage) FindChat(int64) (*storage.Chat, error) {
	return &storage.Chat{}, nil
}
func (t *TestStorage) GetTriggerWords() (*[]storage.TriggerWord, error) {
	return &[]storage.TriggerWord{}, nil
}
func (t *TestStorage) Report(chatID, userID int64, severity int, action, comment string) error {
	return nil
}

type TestBot struct {}

func (t *TestBot) Send(m tgbotapi.Chattable) (tgbotapi.Message, error) {
	return tgbotapi.Message{}, nil
}
func (t *TestBot) AddCommand(c Command) {

}
func (t *TestBot) HandleUpdate(update tgbotapi.Update) []error {
	return nil
}

func (t *TestBot) Start() {}

func (t *TestBot) BanUser(chatID int64, userID int64, revokeMessages bool) error {
	return nil
}
func (t *TestBot) GetCommands() []Command {
	return []Command{}
}

func (t *TestBot) GetUserBio(*tgbotapi.User) (string, error) {
	return "", nil
}

func (t *TestBot) SendToAdmin(m string) error {
	return nil
}
