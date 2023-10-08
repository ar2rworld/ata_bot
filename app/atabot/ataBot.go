package atabot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/ar2rworld/ata_bot/app/commands"
	"github.com/ar2rworld/ata_bot/app/myerror"
	"github.com/ar2rworld/ata_bot/app/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AtaBot struct {
	Bot *tgbotapi.BotAPI
	StartMessage string
	Storage *storage.AtaStorage
	AdminID int64
	updatesChannel *tgbotapi.UpdatesChannel
	commands []commands.Command
}

func NewAtaBot(bot *tgbotapi.BotAPI, startMessage string, storage storage.AtaStorage, updatesChannel *tgbotapi.UpdatesChannel, adminID int64) *AtaBot {
	return &AtaBot{
		Bot: bot,
		StartMessage: startMessage,
		Storage: &storage,
		updatesChannel: updatesChannel,
		AdminID: adminID,
	}
}

func (b *AtaBot) Send(m tgbotapi.Chattable) (tgbotapi.Message, error) {
	return b.Bot.Send(m)
}

func (b *AtaBot) AddCommand(c commands.Command) {
	c.SetAtaBot(b)
	c.SetStorage(*b.Storage)
	c.SetAuthorisedID(b.AdminID)
	b.commands = append(b.commands, c)
}

func (b *AtaBot) HandleUpdate(update tgbotapi.Update) []error {
	var errors []error
	for _, c := range b.commands {
		err := c.Exec(&update)
		if err != nil {
			errors = append(errors,  myerror.NewError(fmt.Sprintf("%s: %s", c.GetName(),err.Error())))
		}
	}

	return errors
}

func (b *AtaBot) Start() {
	for update := range *b.updatesChannel {
		errors := b.HandleUpdate(update)
		if len(errors) > 0 {
			for _, err := range errors {
				if err.Error() == myerror.BadRequestMessageCantBeDeleted {
					message := tgbotapi.NewMessage(update.Message.Chat.ID, "Cannot delete message")
					_, err = b.Send(message)
				}
				log.Println(err)
			}
		}
	}
}

func (b *AtaBot) BanUser(chatID int64, userID int64, revokeMessages bool) error {
	chatMemberConfig := &tgbotapi.ChatMemberConfig{}
	chatMemberConfig.ChatID = chatID
	chatMemberConfig.UserID = userID

	banMember := &tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: *chatMemberConfig,
		RevokeMessages: revokeMessages,
	}
	_, err := b.Send(banMember)
	return err
}

func (b *AtaBot) GetCommands() []commands.Command {
	return b.commands
}

func (b *AtaBot) GetUserBio(u *tgbotapi.User) (string, error) {
	if u.UserName == "" {
		return "", myerror.NewError("user missing username: " + u.String())
	}

	// request telegram page and parse possible bio
	res, err := http.Get(fmt.Sprintf("https://t.me/%s", u.UserName))
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	page := string(data)

	// page contains div with class "tgme_page_description"
	re := regexp.MustCompile("<div\\s*class=\"tgme_page_description\\s*\">(.*)<\\/")
	matches := re.FindStringSubmatch(page)
	if len(matches) < 2 {
		return "", nil
	}
	return matches[1], nil
}