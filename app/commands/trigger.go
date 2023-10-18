package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ar2rworld/ata_bot/app/api"
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
			bio = strings.ToLower(bio)
			
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

			chatID := update.Message.Chat.ID

			switch tempSeverity {
				case storage.Severity200:
					banned, err := ataStorage.IsBanned(&newMember)
					if err != nil {
						return err
					}
					if ! banned {
						err = ataBot.BanUser(chatID, newMember.ID, true)
						if err != nil {
							return err
						}
					}
					err = ataStorage.AddToBanned(&newMember)
					if err != nil {
						return err
					}

					err = ataBot.DeleteMessage(chatID, update.Message.MessageID)
					if err != nil {
						return err
					}

					err = ataStorage.Report(chatID, newMember.ID, storage.Severity200, storage.ActionBanned, "bio:" + triggeredWord)
					if err != nil {
						return err
					}
				break
				// case storage.Severity150:
					// mb just mute?
				// break
				case storage.Severity100:
					err = ataStorage.Report(chatID, newMember.ID, storage.Severity100, storage.ActionNotified, "sus bio:" + triggeredWord)
					if err != nil {
						return err
					}

					messageText := fmt.Sprintf(`suspicious user(%d) bio: "%s" in chat(%d)`, newMember.ID, triggeredWord, chatID)

					data := fmt.Sprintf("%s|,|%d|,|%d|,|%d", BAN, chatID, newMember.ID, update.Message.MessageID)
					susButton := tgbotapi.NewInlineKeyboardButtonData("ban user", data)
					urlButton := tgbotapi.NewInlineKeyboardButtonURL("profile", "https://t.me/" + newMember.UserName)
					markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{ susButton, urlButton })

					m := tgbotapi.NewMessage(ataBot.GetAdminID(), messageText)
					m.ReplyMarkup = markup
					_, err := ataBot.Send(m)
					if err != nil {
						return err
					}
				break
			}
			
			// AnalyzeUserPic if user's bio does not have suspicious words
			// but could have unsafe profile pic
			if tempSeverity != 0 {
				return nil
			}

			apires, err := ataBot.AnalyzeUserPic(&newMember)
			if err != nil {
				return err
			}

			unsafe := apires.Unsafe

			if ! unsafe {
				for _, obj := range apires.Objects {
					if (
						( obj.Label == api.EXPOSED_ANUS ||
						obj.Label == api.EXPOSED_BELLY ||
						obj.Label == api.EXPOSED_BREAST_F ||
						obj.Label == api.EXPOSED_BREAST_M ||
						obj.Label == api.EXPOSED_BUTTOCKS ||
						obj.Label == api.EXPOSED_FEET ||
						obj.Label == api.EXPOSED_GENITALIA_F ||
						obj.Label == api.EXPOSED_GENITALIA_M ) &&
						obj.Score >= 0.5 ) {
						unsafe = true
					}
				}
			}

			if unsafe {
				err = ataStorage.Report(chatID, newMember.ID, storage.Severity100, storage.ActionNotified, "unsafe profile pic")
				if err != nil {
					return err
				}

				messageText := fmt.Sprintf(`unsafe user(%d) profile pic: "%s" in chat(%d)`, newMember.ID, triggeredWord, chatID)

				data := fmt.Sprintf("%s|,|%d|,|%d|,|%d", BAN, chatID, newMember.ID, update.Message.MessageID)
				susButton := tgbotapi.NewInlineKeyboardButtonData("ban user", data)
				urlButton := tgbotapi.NewInlineKeyboardButtonURL("profile", "https://t.me/" + newMember.UserName)
				markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{ susButton, urlButton })

				m := tgbotapi.NewMessage(ataBot.GetAdminID(), messageText)
				m.ReplyMarkup = markup
				_, err := ataBot.Send(m)
				if err != nil {
					return err
				}
			}
			
			// send analysis data to admin
			resData, err := json.Marshal(apires)
			if err != nil {
				return err
			}
			report := fmt.Sprintf("user: %s with id: %d\ndata: %d", newMember.UserName, newMember.ID, resData)
			err = ataBot.SendToAdmin(report)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
