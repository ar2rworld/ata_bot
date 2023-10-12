package storage

import (
	"context"

	"github.com/ar2rworld/ata_bot/app/myerror"
	"go.mongodb.org/mongo-driver/bson"
)

const triggerwords = "triggerwords"

const Severity100 = 100
const Severity150 = 150
const Severity200 = 200

type TriggerWord struct {
	Text string `json:"text" bson:"text,omitempty"`
	Severity int `json:"severity" bson:"severity,omitempty"`
}

func (s *Storage) GetTriggerWords() (*[]TriggerWord, error) {
	voidTriggerWords := &[]TriggerWord{}

	c, err := s.DB.Collection(triggerwords).Find(context.TODO(), bson.D{})
	if err != nil && err.Error() != myerror.NoDocuments {
		return voidTriggerWords, err
	}
	defer c.Close(context.TODO())

	results := []TriggerWord{}
	if err = c.All(context.TODO(), &results); err != nil {
		return voidTriggerWords, err
	}
	
	s.triggerWords = &results

	return s.triggerWords, nil
}
