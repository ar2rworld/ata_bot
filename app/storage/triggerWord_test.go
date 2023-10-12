package storage

import (
	"testing"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetTriggerWords(t *testing.T) {
	sth := NewStorageTestHelper(t)

	mt := sth.MongoTest
	defer mt.Close()

	mt.Run("single GetTriggerWords", func(mt *mtest.T) {
		// targetWord     := "asdf"
		// targetSeverity := Severity100

		// sth.AddMockResponse(mt,
		// 	bson.D{
		// 		{Key: "ok", Value: 1},
		// 		// {Key: "found", Value: bson.A{
		// 			// bson.D{
		// 				{Key: "text", Value: targetWord},
		// 				{Key: "severity", Value: targetSeverity},
		// 			// },
		// 		// },},
		// 	},
		// )
		// // {Key: "triggerwords", Value: bson.A{
		// // 	bson.D{
		// // 		{Key: "text", Value: targetWord},
		// // 		{Key: "severity", Value: targetSeverity},
		// // 	},
		// // },},

		// s := sth.CreateTestStorage(mt)
		// // triggerWords,
		// err := s.TryTriggerWords()

		// if err != nil {
		// 	t.Fatal(err)
		// }
		// // if triggerWords == nil {
		// // 	t.Fatal("TriggerWords is nil")
		// // }
		// // if len(*triggerWords) != 1 {
		// // 	t.Fatalf("Invalid number of trigger words returned: %d", len(*triggerWords))
		// // }

		// // tw := (*triggerWords)[0]

		// // if tw.Text != targetWord {
		// // 	t.Errorf("Incorrect text: \"%s\" != \"%s\"", tw.Text, targetWord)
		// // }
		// // if tw.Severity != targetSeverity {
		// // 	t.Errorf("Incorrect severity: %d != %d", tw.Severity, targetSeverity)
		// // }
	})
}