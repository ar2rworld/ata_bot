package commands

import (
	"github.com/ar2rworld/ata_bot/app/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AtaBotInterface interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	AddCommand(Command)
	GetCommands() []Command
	HandleUpdate(tgbotapi.Update) []error
	Start()
	BanUser(chatID int64, userID int64, revokeMessages bool) error
	GetUserBio(*tgbotapi.User) (string, error)
	SendToAdmin(string) error
	DeleteMessage(chatID int64, messageID int) error
}

type Command interface {
	GetName() string
	Exec(*tgbotapi.Update) error
	GetAtaBot() *AtaBotInterface
	SetAtaBot(AtaBotInterface)
	GetStorage() *storage.AtaStorage
	SetStorage(storage.AtaStorage)
	Authorised(int64) bool
	SetAuthorisedID(int64)
}

type CommandStruct struct {
	adminID int64
	storage *storage.AtaStorage
	ataBot AtaBotInterface
}

func (c *CommandStruct) GetAtaBot() *AtaBotInterface {
	return &c.ataBot
}
func (c *CommandStruct) SetAtaBot(b AtaBotInterface) {
	c.ataBot = b
}
func (c *CommandStruct) GetStorage() *storage.AtaStorage {
	return c.storage
}
func (c *CommandStruct) SetStorage(s storage.AtaStorage) {
	c.storage = &s
}

func (c *CommandStruct) Authorised(id int64) bool {
	if id == c.adminID {
		return true
	}
	return false
}
func (c *CommandStruct) SetAuthorisedID(id int64) {
	c.adminID = id
}