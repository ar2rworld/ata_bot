package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type WatchReply struct {
	CommandStruct
}

func NewWatchReply() *WatchReply {
	return &WatchReply{}
}

func (n *WatchReply) GetName() string {
	return "WatchReply"
}
func (n *WatchReply) GetHelp() string {
	return "watch for forwards from unwanted/unsafe/spam channels or groups"
}

func (c *WatchReply) Exec(u *tgbotapi.Update) error {
	if u.Message == nil {
		return nil
	}

	ataBot     := *c.GetAtaBot()
	ataStorage := *c.GetStorage()

	// check if the message was sent via banned bot
	if u.Message.ViaBot != nil {
		bot := u.Message.ViaBot
		banned, err := ataStorage.IsBanned(bot)
		if err != nil {
			return err
		}
		if banned {
			err := ataBot.BanUser(u.Message.Chat.ID, bot.ID, true)
			if err != nil {
				return err
			}
		}
	}

	// check if the message was forwarded from banned user or channel

	return nil
}
