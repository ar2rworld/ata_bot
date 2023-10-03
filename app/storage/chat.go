package storage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type Chat struct {
	tgbotapi.Chat
	IsActive bool
	Status   string
}

func NewChat(c *tgbotapi.Chat) *Chat {
	t := &Chat{}
	t.ID = c.ID
	t.Type = c.Type
	t.Title = c.Title
	t.UserName = c.UserName
	t.FirstName = c.FirstName
	t.LastName = c.LastName
	t.Bio = c.Bio
	t.HasPrivateForwards = c.HasPrivateForwards
	t.Description = c.Description
	t.InviteLink = c.InviteLink
	t.SlowModeDelay = c.SlowModeDelay
	t.MessageAutoDeleteTime = c.MessageAutoDeleteTime
	t.HasProtectedContent = c.HasProtectedContent
	t.StickerSetName = c.StickerSetName
	t.CanSetStickerSet = c.CanSetStickerSet
	t.LinkedChatID = c.LinkedChatID
	return t
}

func (c *Chat) toDoc () bson.D {
	return bson.D{
		{ Key: "id", Value: c.ID },
		{ Key: "type", Value: c.Type },
		{ Key: "title", Value: c.Title },
		{ Key: "username", Value: c.UserName },
		{ Key: "firstname", Value: c.FirstName },
		{ Key: "lastname", Value: c.LastName },
		{ Key: "hasprivateforwards", Value: c.HasPrivateForwards },
		{ Key: "description", Value: c.Description },
		{ Key: "slowmodedelay", Value: c.SlowModeDelay },
		{ Key: "messageautodeletetime", Value: c.MessageAutoDeleteTime },
		{ Key: "hasprotectedcontent", Value: c.HasProtectedContent },
		{ Key: "stickersetname", Value: c.StickerSetName },
		{ Key: "cansetstickerset", Value: c.CanSetStickerSet },
		{ Key: "linkedchatid", Value: c.LinkedChatID },
		{ Key: "isactive", Value: c.IsActive },
		{ Key: "status", Value: c.Status },
	}
}
