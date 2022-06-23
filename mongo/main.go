package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ext-tg-bot/utils"
)

var (
	_mongoClient *mongo.Client
	_mongoDB     *mongo.Database
	_ctx         context.Context
)

// Init initializes Mongo DB connection
func Init(ctx context.Context, mongoURL string, mongoDB string) error {
	var err error
	_ctx = ctx

	utils.LogInfo("Starting Mongo:", mongoURL)

	_mongoClient, err = mongo.Connect(_ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return err
	}
	_mongoDB = _mongoClient.Database(mongoDB)

	return nil
}

// Close closes Mongo DB connections
func Close() {
	err := _mongoClient.Disconnect(_ctx)
	if err != nil {
		utils.LogError(err)
	}
}
