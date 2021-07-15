package mongo

// keep these organized by category with an empty line between each
// 1. core
// 2. remote
// 3. local
import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Cx *MongoDatastore

type MongoDatastore struct {
	db      *mongo.Database
	Session *mongo.Client
}

var (
	ConnectToMongo = connectToMongo
)

func init() {
	addr := os.Getenv("MONGO_URL")
	database := os.Getenv("MONGO_DATABASE")

	var mongoDataStore *MongoDatastore
	db, session := connect(addr, database)
	if db != nil && session != nil {
		mongoDataStore = new(MongoDatastore)
		mongoDataStore.db = db
		mongoDataStore.Session = session
		Cx = mongoDataStore
	}
}

func connect(addr string, database string) (a *mongo.Database, b *mongo.Client) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	connectOnce.Do(func() {
		session = ConnectToMongo(addr, database)
		if session != nil {
			db = session.Database(database)
		}
	})

	if session != nil {
		return db, session
	}

	return nil, nil
}

func connectToMongo(addr string, database string) (b *mongo.Client) {
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(addr)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil
	}

	return client
}
