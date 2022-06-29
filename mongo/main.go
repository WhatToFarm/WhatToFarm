package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ton-tg-bot/logger"
)

const (
	userCollection = "user"
)

var (
	_mongoClient *mongo.Client
	_mongoDB     *mongo.Database
	_ctx         context.Context
)

// Init initializes Mongo DB connection
func Init(ctx context.Context, mongoURL, mongoDB string) error {
	var err error
	_ctx = ctx

	logger.LogInfo("Starting Mongo:", mongoURL)

	_mongoClient, err = mongo.Connect(_ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return fmt.Errorf("mongoDB connect: %w", err)
	}
	_mongoDB = _mongoClient.Database(mongoDB)

	if err = createIndex(); err != nil {
		return err
	}

	return nil
}

// Close closes Mongo DB connections
func Close() {
	err := _mongoClient.Disconnect(_ctx)
	if err != nil {
		logger.LogError(err)
	}
}
