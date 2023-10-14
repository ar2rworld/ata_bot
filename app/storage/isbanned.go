package storage

import (
	"context"

	"github.com/ar2rworld/ata_bot/app/myerror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Storage) IsBanned(user *tgbotapi.User) (bool, error) {
	var probablyFake FakeUser
	err := s.DB.Collection(FakesCollection).FindOne(
		context.TODO(),
		bson.D{{Key: "id", Value: user.ID}},
	).Decode(&probablyFake)

	if err != nil && err.Error() == myerror.NoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return probablyFake.IsBanned, nil
}
