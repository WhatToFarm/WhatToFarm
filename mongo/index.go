package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"ton-tg-bot/models"
)

func createIndex() error {
	uniqTrue := true

	coll := _mongoDB.Collection(userCollection)
	index := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: models.FieldAccount, Value: 1},
			},
			Options: &options.IndexOptions{
				Unique: &uniqTrue,
			},
		},
	}
	opts := options.CreateIndexes().SetMaxTime(30 * time.Second)
	if _, err := coll.Indexes().CreateMany(_ctx, index, opts); err != nil {
		return fmt.Errorf("index: %w", err)
	}

	return nil
}
