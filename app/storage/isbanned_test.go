package storage

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestIsBanned(t *testing.T) {
	sth := NewStorageTestHelper(t)
	mt := sth.MongoTest
	defer mt.Close()

	mt.Run("IsBanned no user found", func(mt *mtest.T) {
		sth.AddMockResponseNoDocumentsFound(mt)

		s := sth.CreateTestStorage(mt)

		testUserID := int64(1)
		user := &tgbotapi.User{ ID: testUserID}

		banned, err := s.IsBanned(user)
		if err != nil {
			t.Fatalf("Should not return error, but got: %s", err)
		}

		if banned {
			t.Error("Non existing user should not be banned")
		}
	})

	mt.Run("IsBanned user found", func(mt *mtest.T) {
		testUserID       := int64(1)
		testUserIsBanned := true

		sth.AddMockResponse(mt, bson.D{
			{Key: "id", Value: testUserID},
			{Key: "isbanned", Value: testUserIsBanned},
		})

		s := sth.CreateTestStorage(mt)

		user := &tgbotapi.User{ ID: testUserID}

		banned, err := s.IsBanned(user)
		if err != nil {
			t.Fatalf("Should not return error, but got: %s", err)
		}

		if ! banned {
			t.Error("User found and should be banned")
		}
	})
}