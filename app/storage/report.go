package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

const reports = "reports"

type ReportStruct struct {
	ChatID 	 int64
	UserID 	 int64
	Severity int
	Action   string
	Comment  string
}
func NewReport(chatID, userID int64, severity int, action, comment string) *ReportStruct {
	return &ReportStruct{
		ChatID: chatID,
		UserID: userID,
		Severity: severity,
		Action: action,
		Comment: comment,
	}
}
func (r *ReportStruct) toDoc() bson.D {
	return bson.D{
		{Key: "chatid", Value: r.ChatID},
		{Key: "userid", Value: r.UserID},
		{Key: "severity", Value: r.Severity},
		{Key: "action", Value: r.Action},
		{Key: "comment", Value: r.Comment},
	}
}

func (s *Storage) Report(chatID, userID int64, severity int, action, comment string) error {
	r := NewReport(chatID, userID, severity, action, comment)

	_, err := s.DB.Collection(reports).InsertOne(context.TODO(), r.toDoc())
	if err != nil {
		return err
	}
	
	return nil
}