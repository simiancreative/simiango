package mongo

import (
	"errors"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestConnect(t *testing.T) {
	mockError := errors.New("uh oh")
	testing.Init()
	subtests := []struct {
		name           string
		u              string
		ConnectHandler func(string, string) (*mongo.Database, *mongo.Client, error)
		ConnectToMongo func(string, string) (*mongo.Client, error)
		expectedErr    error
	}{
		{
			name: "Test Connect",
			ConnectToMongo: func(s string, s2 string) (*mongo.Client, error) {
				if s != "//u:p@a/db" {
					return nil, errors.New("wrong connection string")
				}
				return nil, nil
			},
			ConnectHandler: func(s string, s2 string) (a *mongo.Database, b *mongo.Client, c error) {
				b, c = ConnectToMongo(s, s2)
				return nil, b, c
			},
		},
		{
			name: "error from Connect",
			ConnectToMongo: func(s string, s2 string) (*mongo.Client, error) {
				return nil, mockError
			},
			ConnectHandler: func(s string, s2 string) (a *mongo.Database, b *mongo.Client, c error) {
				b, c = ConnectToMongo(s, s2)
				return nil, b, c
			},
			expectedErr: mockError,
		},
	}

	//addr := os.Getenv("MONGO_URL")
	//database := os.Getenv("MONGO_DATABASE")
	addr := "//u:p@a/db"
	database := "db"

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ConnectHandler = subtest.ConnectHandler
			ConnectToMongo = subtest.ConnectToMongo
			var err error
			fmt.Println(addr)
			_, _, err = ConnectHandler(addr, database)

			if !errors.Is(err, subtest.expectedErr) {
				t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
			}
		})
	}
}
