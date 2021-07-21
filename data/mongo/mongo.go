package mongo

// keep these organized by category with an empty line between each
// 1. core
// 2. remote
// 3. local
import (
	"context"
	"fmt"
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

func connectHandler(addr string, database string) (a *mongo.Database, b *mongo.Client, c error) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	var err error
	fmt.Println("connecting...")
	connectOnce.Do(func() {
		session, err := connectToMongo(addr, database)
		if err == nil {
			db = session.Database(database)
		}
	})

	if session != nil {
		return db, session, nil
	}

	return nil, nil, err
}

func connectToMongo(addr string, database string) (b *mongo.Client, c error) {
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(addr)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
