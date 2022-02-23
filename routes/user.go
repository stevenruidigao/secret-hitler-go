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

func RegisterUser(user *types.UserPrivate) *types.UserPrivate {
	if user == nil {
		return nil
	}

	user.UserPublic.Username = strings.ReplaceAll(user.UserPublic.Username, " ", "-")

	if user.UserPublic.Username == "" {
		return nil
	}

	cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
		"userPublic.username": user.UserPublic.Username,
	})

	ok := true

	for cursor.Err() != mongo.ErrNoDocuments {
		ok = false
		user.UserPublic.Username = "*" + user.UserPublic.Username

		cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.username": user.UserPublic.Username,
		})
	}

	userID := uuid.NewString()

	cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
		"userPublic.ID": userID,
	})

	for cursor.Err() != mongo.ErrNoDocuments {
		userID = uuid.NewString()

		cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.ID": userID,
		})
	}

	user.UserPublic.ID = userID
	user.UserPublic.Created = time.Now()
	user.UserPublic.EloOverall = 1600
	user.UserPublic.EloSeason = 1600
	user.UserPublic.Status = nil
	user.UserPublic.Profile.UserID = userID
	user.UserPublic.Profile.Username = user.UserPublic.Username
	user.UserPublic.Profile.Created = time.Now()
	user.UserPublic.Profile.LastConnected = time.Now()
	user.UserPublic.Profile.Badges = []types.Badge{}
	user.UserPublic.Profile.RecentGames = []types.RecentGame{}
	user.GameSettings.Blacklist = []string{}
	user.FinishedSignup = ok
	database.MongoDB.Collection("Users").InsertOne(ctx, user)
	return user
}

func AddSessionToUser(session *types.Session, user *types.UserPrivate) {
	session.UserID = user.UserPublic.ID
	session.Expires = time.Now().Add(7 * 24 * time.Hour)
	database.MongoDB.Collection("Sessions").InsertOne(ctx, *session)
	user.Sessions = append(user.Sessions, *session)
	database.UpdateUserByID(user.UserPublic.ID, user)
}

func UpdateUserUsername(user *types.UserPrivate, username string) bool {
	cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
		"userPublic.username": username,
	})

	if cursor.Err() != mongo.ErrNoDocuments {
		return false
	}

	user.UserPublic.Username = username
	database.UpdateUserByID(user.UserPublic.ID, user)

	return true
}
