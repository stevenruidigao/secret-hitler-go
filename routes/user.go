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

	ok := true

	for cursor.Err() != mongo.ErrNoDocuments {
		ok = false
		username = "*" + username

		cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.username": username,
		})
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
		FinishedSignup: ok,
	}

	database.MongoDB.Collection("Users").InsertOne(ctx, user)
	return &user
}

func AddSessionToUser(session *types.Session, user *types.UserPrivate) {
	session.UserID = user.UserPublic.UserID
	session.Expires = time.Now().Add(7 * 24 * time.Hour)
	database.MongoDB.Collection("Sessions").InsertOne(ctx, *session)
	user.Sessions = append(user.Sessions, *session)
}
