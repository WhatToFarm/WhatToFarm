package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"ton-tg-bot/logger"
	"ton-tg-bot/models"
)

func CreateUser(user *models.TgUser) error {
	_, err := _mongoDB.Collection(userCollection).InsertOne(_ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func GetUsers() ([]*models.TgUser, error) {
	cur, err := _mongoDB.Collection(userCollection).Find(_ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}
	defer func() {
		if errClose := cur.Close(_ctx); errClose != nil {
			logger.LogError("cursor close:", errClose)
		}
	}()

	users := make([]*models.TgUser, 0)
	if err = cur.All(_ctx, &users); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}
	return users, nil
}

func UpdateUser(user *models.TgUser) error {
	filter := bson.M{
		models.FieldAccount: user.GitAccount,
	}
	update := bson.M{
		"$set": bson.M{
			models.FieldTS:       user.TS,
			models.FieldAttempts: user.Attempts,
		},
	}
	_, err := _mongoDB.Collection(userCollection).UpdateOne(_ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
