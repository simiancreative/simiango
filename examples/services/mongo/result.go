package mongoexample

import (
	"context"
	"log"

	m "github.com/simiancreative/simiango/data/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (s Service) Result() (interface{}, error) {
	rows := []Product{}

	var ctx = context.Background()
	collection := m.Cx.Session.Database("ExampleData").Collection("products")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var result Product
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		rows = append(rows, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	return rows, err
}
