package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/ar2rworld/ata_bot/app/atabot"
	"github.com/ar2rworld/ata_bot/app/commands"
	"github.com/ar2rworld/ata_bot/app/storage"
	"go.mongodb.org/mongo-driver/bson"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	adminID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	ataStorage, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}

	unbanArt(bot, ataStorage)

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	ataBot := atabot.NewAtaBot(
		bot,
		"Hello, son",
		ataStorage,
		&updates,
		adminID,
	)

	ataBot.AddCommand(commands.NewMaskara())
	ataBot.AddCommand(commands.NewNewGroupUpdate())
	ataBot.AddCommand(commands.NewNewGroupMember())
	ataBot.AddCommand(commands.NewTrigger())
	ataBot.AddCommand(commands.NewHelp())

	err = ataBot.SendToAdmin("Bastaimyn goi")
	if err != nil {
		log.Println("Start message error:", err)
	}

	ataBot.Start()
}

func unbanArt(bot *tgbotapi.BotAPI, ataStorage *storage.Storage){
	artID  := int64(1265820975)
	chatID := int64(-1001506079405)
	chatMemberConfig := &tgbotapi.ChatMemberConfig{
		ChatID: chatID,
		UserID: artID,
	}

	message := &tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: *chatMemberConfig,
	}

	_, err := bot.Send(message)
	if err != nil && err.Error() != "json: cannot unmarshal bool into Go value of type tgbotapi.Message" {
		log.Println("unban:", err)
	}
	result, err := ataStorage.DB.Collection(storage.FakesCollection).DeleteOne(context.Background(), bson.D{{Key: "id", Value: artID}})
	if err != nil {
		log.Println(err)
	}
	log.Println("deleted: ", result.DeletedCount)
}