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
		"userPublic.userID": userID,
	})

	if cursor.Err() != nil {
		return nil
	}

	user := types.UserPrivate{}
	cursor.Decode(&user)

	return &user
}

func UpdateUserByID(userID string, user *types.UserPrivate) bool {
	if user == nil {
		return false
	}

	MongoDB.Collection("Users").UpdateOne(ctx, bson.M{
		"userPublic.userID": userID,
	}, bson.M{
		"$set": user,
	})

	return true
}
