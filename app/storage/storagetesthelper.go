package storage

import (
	"testing"

	"github.com/ar2rworld/ata_bot/app/myerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type StorageTestHelper struct {
	MongoTest *mtest.T
	Client *mongo.Client
	DB *mongo.Database
	Storage Storage
}

func NewStorageTestHelper(t *testing.T) *StorageTestHelper {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	return &StorageTestHelper{
		MongoTest: mt,
		Client: mt.Client,
		DB: mt.DB,
	}
}
func (sth *StorageTestHelper) CreateTestStorage(mt *mtest.T) *Storage {
	st := &Storage{ Client: *mt.Client, DB: *mt.DB}
	sth.Storage = *st
	return st
}

func (sth *StorageTestHelper) AddMockResponse(mt *mtest.T, res bson.D) {
	mt.AddMockResponses(
		mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, res),
	)
}

func (sth *StorageTestHelper) AddMockResponseNoDocumentsFound(mt *mtest.T) {
	mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
		Code:    11000,
		Message: myerror.NoDocuments,
 }))
}
