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

	// unbanArt(bot, ataStorage)
	// banArt(bot)

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
	ataBot.AddCommand(commands.NewTriggerCallbackQuery())

	err = ataBot.SendToAdmin("Bastaimyn goi")
	if err != nil {
		log.Println("Start message error:", err)
	}

	ataBot.Start()
}

func unbanArt(bot *tgbotapi.BotAPI, ataStorage *storage.Storage) {
	artID := int64(1265820975)
	chatID := int64(-1001506079405)
	chatMemberConfig := &tgbotapi.ChatMemberConfig{
		ChatID: chatID,
		UserID: artID,
	}

	un := &tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: *chatMemberConfig,
	}

	res, err := bot.Request(un)
	if err != nil {
		log.Println("unban:", err)
	}
	if res.ErrorCode != 0 {
		log.Printf("error during unban call: errorCode: %d, description: %s, json raw: %s", res.ErrorCode, res.Description, res.Result)
	}
	result, err := ataStorage.DB.Collection(storage.FakesCollection).DeleteOne(context.Background(), bson.D{{Key: "id", Value: artID}})
	if err != nil {
		log.Println(err)
	}
	log.Println("deleted: ", result.DeletedCount)
}

type BanRequest struct {
	UserID         int64
	ChatID         int64
	revokeMessages bool
}

func (BanRequest) method() string {
	return "banChatMember"
}

func (r BanRequest) params() (tgbotapi.Params, error) {
	return map[string]string{
		"user_id":         strconv.FormatInt(r.UserID, 10),
		"chat_id":         strconv.FormatInt(r.ChatID, 10),
		"revoke_messages": strconv.FormatBool(true),
	}, nil
}

func banArt(bot *tgbotapi.BotAPI) {
	// r := &BanRequest{
	// 	UserID: int64(1265820975),
	// 	ChatID: int64(-1001506079405),
	// }
	// var t tgbotapi.Chattable = r
	// t = BanRequest{
	// 	UserID: int64(1265820975),
	// 	ChatID: int64(-1001506079405),
	// }
	artID := int64(1265820975)
	chatID := int64(-1001506079405)
	r := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: artID,
		},
		RevokeMessages: true,
	}

	res, err := bot.Request(r)
	if err != nil {
		log.Println("err:", err)
	}
	log.Printf("err code: %d, desc: %s, json raw: %s", res.ErrorCode, res.Description, res.Result)
}
