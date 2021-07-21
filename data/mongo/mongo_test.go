package mongo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func init() {
	addr := os.Getenv("MONGO_URL")
	database := os.Getenv("MONGO_DATABASE")
	var mongoDataStore *MongoDatastore

	db, session, _ := connectHandler(addr, database)
	if db != nil && session != nil {
		mongoDataStore = new(MongoDatastore)
		mongoDataStore.db = db
		mongoDataStore.Session = session
		Cx = mongoDataStore
	}
}

func TestConnect(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	ConnectToMongo = func(s string, s2 string) (*mongo.Client, error) {
		return mt.Client, nil
	}

	mt.Run("connect with database name success", func(mt *mtest.T) {
		database := os.Getenv("MONGO_DATABASE")
		//mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

		testing.Init()
		assert.Equal(t, database, Cx.db.Name())
	})
}
