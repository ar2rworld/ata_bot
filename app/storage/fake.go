package storage

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type FakeUser struct {
	tgbotapi.User
	IsBanned bool
}

func (f *FakeUser) toDoc() bson.D {
	return bson.D{
		{ Key: "id", Value: f.ID},
		{ Key: "isbanned", Value: f.IsBanned },
		{ Key: "isbot", Value: f.IsBot },
		{ Key: "firstname", Value: f.FirstName },
		{ Key: "lastname", Value: f.LastName },
		{ Key: "username", Value: f.UserName },
		{ Key: "languagecode", Value: f.LanguageCode },
		{ Key: "canjoingroups", Value: f.CanJoinGroups },
		{ Key: "canrealallgroupmessages", Value: f.CanReadAllGroupMessages },
		{ Key: "supportsinlinequeries", Value: f.SupportsInlineQueries },		
	}
}