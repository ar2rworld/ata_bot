package commands

import (
	"strings"

	"github.com/ar2rworld/ata_bot/app/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Trigger struct {
	CommandStruct
}

func NewTrigger() *Trigger {
	return &Trigger{}
}

func (n *Trigger) GetName() string {
	return "Trigger"
}
func (n *Trigger) GetHelp() string {
	return "when new member joins the chat, bot will check if user is bot or has some trigger words when ban him"
}

func (t *Trigger) Exec(update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}
	if len(update.Message.NewChatMembers) > 0 {
		ataBot  := *t.GetAtaBot()
		ataStorage := *t.GetStorage()
		triggerWords, err := ataStorage.GetTriggerWords()
		if err != nil {
			return err
		}

		for _, newMember := range update.Message.NewChatMembers {
			
			// if IsBot
			if newMember.IsBot {
				err := ataBot.BanUser(update.Message.Chat.ID, newMember.ID, true)
				if err != nil {
					return err
				}
				err = ataStorage.AddToBanned(&newMember)
				if err != nil {
					return nil
				}
			}

			// TODO: if newMember has some interesting words in bio
			if newMember.UserName == "" {
				return nil
			}

			bio, err := ataBot.GetUserBio(&newMember)
			if err != nil {
				return err
			}
			tempSeverity := 0
			triggeredWord := "" 
			for _, tw := range *triggerWords {
				foundTriggerWord := strings.Contains(bio, tw.Text)
				if foundTriggerWord {
					if tempSeverity < tw.Severity {
						triggeredWord = bio
						tempSeverity = tw.Severity
					}
				}
			}
			if tempSeverity == 0 {
				return nil
			}

			chatID := update.Message.Chat.ID

			switch tempSeverity {
				case storage.Severity200:
					err = ataBot.BanUser(chatID, newMember.ID, true)
					if err != nil {
						return err
					}
					err = ataStorage.AddToBanned(&newMember)
					if err != nil {
						return err
					}
					err = ataStorage.Report(chatID, newMember.ID, storage.Severity200, "banned", "bio:" + triggeredWord)
					if err != nil {
						return err
					}
				break
				// case storage.Severity150:
			}
			
		}
	}
	return nil
}
