package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/ar2rworld/ata_bot/app/myerror"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const FakesCollection = "fakes"
const chatsCollection = "chats"

const ActionBanned = "banned"
const ActionNotified = "notified"

type AtaStorage interface {
	IsBanned(*tgbotapi.User) (bool, error)
	AddToBanned(*tgbotapi.User) error
	AddNewChat(*tgbotapi.Chat, string) error
	RemoveChat(*tgbotapi.Chat, string) error
	FindChat(int64) (*Chat, error)
	GetTriggerWords() (*[]TriggerWord, error)
	Report(chatID, userID int64, severity int, action, comment string) error
}

type Storage struct {
	Client        			 mongo.Client
	DB								 	 mongo.Database
	mongoUser     		 	 string
	mongoHost     			 string
	triggerWords				 *[]TriggerWord
}

func NewStorage() (*Storage, error) {
	emptyStorage := &Storage{}
	
	dbName := os.Getenv("DB_NAME")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")

	if dbName == "" ||
	mongoUser == "" ||
	mongoPassword == "" ||
	mongoHost == "" {
		return emptyStorage, myerror.NewError("Couldn't find env vars")
	}

	uri := fmt.Sprintf(`mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority`, mongoUser, mongoPassword, mongoHost)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return emptyStorage, err
	}

	db := client.Database(dbName)
	if err := db.RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return emptyStorage, err
	}
	storage := &Storage{
		Client: *client,
		DB: *db,
		mongoUser: mongoUser,
		mongoHost: mongoHost,
	}
	return storage, nil
}

func (s *Storage) AddToBanned(user *tgbotapi.User) error {
	var fake FakeUser = FakeUser{
		*user,
		true,
	}
	_, err := s.DB.Collection(FakesCollection).InsertOne(
		context.TODO(),
		fake.toDoc(),
	)
	return err
}

func (s *Storage) AddNewChat(c *tgbotapi.Chat, status string) error {
	foundChat, err := s.FindChat(c.ID)

	// no errors means chat is found, no action needed
	if err == nil && ! foundChat.IsActive {
		foundChat.IsActive = true
		foundChat.Status = status
		_, err = s.DB.Collection(chatsCollection).UpdateOne(
			context.TODO(),
			bson.D{{ Key: "id", Value: c.ID }},
			bson.D{{ Key: "$set", Value: foundChat.toDoc() }},
		)
		if err != nil {
			return err
		}
	}

	if err != nil && err.Error() == myerror.NoDocuments {
		newChat := NewChat(c)
		newChat.IsActive = true
		newChat.Status = status

		_, err := s.DB.Collection(chatsCollection).InsertOne(context.TODO(), newChat.toDoc())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) RemoveChat(c *tgbotapi.Chat, status string) error {
	id := c.ID
	chat, err := s.FindChat(id)
	if err != nil {
		return err
	}

	chat.IsActive = false
	chat.Status = status

	_, err = s.DB.Collection(chatsCollection).UpdateOne(
		context.TODO(), bson.D{{ Key: "id", Value: id }},
		bson.D{{
			Key: "$set", Value: chat.toDoc(),
		}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindChat(id int64) (*Chat, error) {
	var emptyChat = &Chat{}

	var tempChat struct {
		IsActive bool
		Status string
	}
	var foundChat tgbotapi.Chat

	result := s.DB.Collection(chatsCollection).FindOne(
		context.TODO(),
		bson.D{
			{ Key: "id", Value: id, },
		},
	)
	
	err := result.Decode(&foundChat)
	if err != nil {
		return emptyChat, err
	}
	err = result.Decode(&tempChat)
	if err != nil {
		return emptyChat, err
	}
	return &Chat{
		Chat: foundChat,
		IsActive: tempChat.IsActive,
		Status: tempChat.Status,		
	}, nil
}
