package routes

import (
	"secrethitler.io/database"
	"secrethitler.io/types"

	"strings"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(username string) *types.UserPrivate {
	cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
		"userPublic.username": username,
	})

	if cursor.Err() != mongo.ErrNoDocuments {
		return nil
	}

	userID := uuid.NewString()

	cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
		"userPublic.userID": userID,
	})

	for cursor.Err() != mongo.ErrNoDocuments {
		userID = uuid.NewString()

		cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.userID": userID,
		})
	}

	user := types.UserPrivate{
		UserPublic: types.UserPublic{
			Username:   strings.ReplaceAll(username, " ", "-"),
			UserID:     userID,
			Created:    time.Now(),
			EloOverall: 1600,
			EloSeason:  1600,
			Status: types.UserStatus{
				Type: "none",
			},
		},
		GameSettings: types.GameSettings{
			Blacklist: []string{},
		},
		Profile: types.Profile{
			Created:     time.Now(),
			RecentGames: []types.GamePublic{},
		},
	}

	database.MongoDB.Collection("Users").InsertOne(ctx, user)
	return &user
}
