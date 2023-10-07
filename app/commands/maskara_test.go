package commands

import (
	"fmt"
	"testing"

	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestMaskara(t *testing.T) {
	adminID := int64(1014210753)
	chatID  := int64(-1001757386839)
	userID  := int64(6384738980)
	
	updateString := fmt.Sprintf(`{"update_id":127503233,"message":{"message_id":841,"from":{"id":1014210753,"is_bot":false,"first_name":"Nemo","last_name":"Cap","username":"ne0ne0postaviat0kolonku","language_code":"en"},
	"chat":{"id":%d,"title":"Name of a group in Kelowna","type":"supergroup"},"date":1696276074,"message_thread_id":835,"reply_to_message":{"message_id":835,
	"from":{"id":%d,"is_bot":false,"first_name":"Latia"},"chat":{"id":-1001757386839,"title":"Name of a group in Kelowna","type":"supergroup"},"date":1696192274,
	"forward_from":{"id":6307729024,"is_bot":true,"first_name":"LocalGirls","username":"LocalGirls69_bot"},"forward_date":1696191909,"animation":{"file_name":"1.mp4","mime_type":"video/mp4","duration":6,"width":404,"height":720,"thumbnail":{"file_id":"someID","file_unique_id":"someID","file_size":10847,"width":180,"height":320},"thumb":{"file_id":"someID","file_unique_id":"someID","file_size":10847,"width":180,"height":320},"file_id":"someID","file_unique_id":"someID","file_size":1775793},"document":{"file_name":"1.mp4","mime_type":"video/mp4","thumbnail":{"file_id":"AAMCBQADHQJov5hXAAIDQ2UbHcLE4u_O-Sob6iJ_fiGDM4vAAAKyCQACmGvQVMHJaVhb26V4AQAHbQADMAQ","file_unique_id":"AQADsgkAAphr0FRy","file_size":10847,"width":180,"height":320},"thumb":{"file_id":"AAMCBQADHQJov5hXAAIDQ2UbHcLE4u_O-Sob6iJ_fiGDM4vAAAKyCQACmGvQVMHJaVhb26V4AQAHbQADMAQ","file_unique_id":"AQADsgkAAphr0FRy","file_size":10847,"width":180,"height":320},"file_id":"CgACAgUAAx0CaL-YVwACA0NlGx3CxOLvzvkqG-oif34hgzOLwAACsgkAAphr0FTByWlYW9uleDAE","file_unique_id":"AgADsgkAAphr0FQ","file_size":1775793},"reply_markup":{"inline_keyboard":[[{"text":"\ud83d\udccc Whatsapp groups-local \ud83d\udccc","url":"someUrl"}],[{"text":"\ud83c\udf6d Telegram groups-local \ud83c\udf6d","url":"someUrl"}],[{"text":"\ud83c\udfa5 live stream hot girls \ud83c\udfa5","url":"someUrl"}]]}},"text":"/maskara","entities":[{"offset":0,"length":8,"type":"bot_command"}]}}`,
		chatID,
		userID,
	)

	update, err := NewUpdate(updateString)
	if err != nil {
		t.Error(err)
	}

	testBot := &maskaraTestBot{}
	testStorage := &TestStorage{}

	m := NewMaskara()
	m.SetAtaBot(testBot)
	m.SetStorage(testStorage)
	m.SetAuthorisedID(adminID)

	err = m.Exec(update)

	if err != nil {
		t.Fatal(err)
	}
	if testBot.BannedChatID != chatID && testBot.BannedUserID != userID {
		t.Fatalf(
			"incorrect banned ids, chatID: %d != %d; userID: %d != %d",
			testBot.BannedChatID, chatID,
			testBot.BannedUserID, userID,
		)
	}

}

type maskaraTestBot struct {
	TestBot
	BannedChatID int64
	BannedUserID int64
}
func (t *maskaraTestBot) BanUser(chatID, userID int64, revokeMessages bool) error {
	t.BannedChatID = chatID
	t.BannedUserID = userID
	return nil
}

// {"update_id":127503290,
//"my_chat_member":{"chat":{"id":-1001506079405,"title":"Public group for my bots","username":"my_bots_group","type":"supergroup"},"from":{"id":1014210753,"is_bot":false,"first_name":"Nemo","last_name":"Cap","username":"ne0ne0postaviat0kolonku","language_code":"en"},"date":1696313736,"old_chat_member":{"user":{"id":6592599799,"is_bot":true,"first_name":"Ata bot","username":"the_ata_bot"},"status":"member"},"new_chat_member":{"user":{"id":6592599799,"is_bot":true,"first_name":"Ata bot","username":"the_ata_bot"},"status":"administrator","can_be_edited":false,"can_manage_chat":true,"can_change_info":true,"can_delete_messages":true,"can_invite_users":true,"can_restrict_members":true,"can_pin_messages":true,"can_manage_topics":false,"can_promote_members":false,"can_manage_video_chats":true,"is_anonymous":false,"can_manage_voice_chats":true,"custom_title":"ata"}}}
