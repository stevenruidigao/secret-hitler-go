package database

import (
	"secrethitler.io/types"

	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()
var MongoDB *mongo.Database
var RedisDB *redis.Client

func Marshal(data interface{}) string {
	result, _ := json.Marshal(data)

	return string(result)
}

func SetupDatabase(mongoDB *mongo.Database, redisDB *redis.Client) {
	MongoDB = mongoDB
	RedisDB = redisDB
	RedisDB.Set(ctx, "playerCount", "0", 0)
}

func GetUserByID(userID string) *types.UserPrivate {
	cursor := MongoDB.Collection("Users").FindOne(ctx, bson.M{
		"user.userID": userID,
	})

	if cursor.Err() != nil {
		return nil
	}

	user := types.UserPrivate{}
	cursor.Decode(&user)

	return &user
}
