package routes

import (
	"secrethitler.io/database"
	"secrethitler.io/socket"
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	// "html/template"
	"net/http"

	//	"strconv"
	"strings"
	"time"

	//	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	// "github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()

func Authenticate(request *http.Request) *types.Session {
	session, _ := Store.Get(request, "sh-session")
	localSession, ok := session.Values["session"].(types.Session)

	if !ok {
		return nil
	}

	// fmt.Println(localSession)

	cursor := database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
		"token":  localSession.Token,
		"userID": localSession.UserID,
	})

	// fmt.Println(cursor.Err(), "???")

	if cursor.Err() != nil {
		return nil
	}

	cursor.Decode(&localSession)

	if time.Until(localSession.Expires) < 0 {
		fmt.Println("out of time")

		return nil
	}

	return &localSession
}

func Marshal(data interface{}) string {
	result, _ := json.Marshal(data)
	return string(result)
}

func SetupRoutes(router *mux.Router, io *socketio.Server, store *sessions.CookieStore) {
	Store = store
	gothic.Store = store
	gob.Register(types.Session{})
	router.Handle("/", Render("page-home"))
	router.Handle("/404", Render("404")).Methods("GET")
	router.Handle("/game/", Render("game")).Methods("GET")
	router.Handle("/game/*", Render("game")).Methods("GET")
	router.Handle("/oauth-select-username", Render("page-new-username")).Methods("GET")

	router.HandleFunc("/online-playercount", func(writer http.ResponseWriter, request *http.Request) {
		data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
		utils.JSONResponse(writer, struct {
			Count string `json:"count"`
		}{data}, 200)
	})

	router.HandleFunc("/profile", func(writer http.ResponseWriter, request *http.Request) {
		// fmt.Println(request.URL.Query()["username"][0])

		cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.username": request.URL.Query()["username"][0],
		})

		if cursor.Err() != nil {
			fmt.Println("*&&&&&", cursor.Err())

			return
		}

		var user types.UserPublic
		cursor.Decode(&user)
		// fmt.Println("*&&&&&", user, user.Profile, ";;;;;;;;;;", user.Profile.RecentGames, len(user.Profile.RecentGames))
		utils.JSONResponse(writer, user.Profile, 200)
	})

	router.Handle("/socket.io/", socket.SetupSocketRoutes(io, store))

	router.HandleFunc("/{provider}-login", func(writer http.ResponseWriter, request *http.Request) {
		_, err := gothic.CompleteUserAuth(writer, request)

		if err == nil {
			writer.Header().Set("Location", "/game/")
			writer.WriteHeader(http.StatusTemporaryRedirect)

		} else {
			gothic.BeginAuthHandler(writer, request)
		}
	})

	router.HandleFunc("/auth/{provider}/callback", func(writer http.ResponseWriter, request *http.Request) {
		user, err := gothic.CompleteUserAuth(writer, request)

		if err != nil {
			fmt.Println(err)
			writer.Header().Set("Location", "/")
			writer.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		localSession := types.Session{
			Token: uuid.NewString(),
		}

		cursor := database.MongoDB.Collection("Sessions").FindOne(ctx, localSession)

		for cursor.Err() != mongo.ErrNoDocuments {
			localSession.Token = uuid.NewString()
			cursor = database.MongoDB.Collection("Sessions").FindOne(ctx, localSession)
		}

		cursor = database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"linkedAccounts.provider": user.Provider,
			"linkedAccounts.userid":   user.UserID,
		})

		localUser := &types.UserPrivate{
			UserPublic: types.UserPublic{
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

		// fmt.Println(cursor.Err(), user, user.NickName, user.Name, user.Email, "**")

		if cursor.Err() == mongo.ErrNoDocuments {
			localUser.UserPublic.Username = user.NickName

			if localUser.Username == "" {
				localUser.UserPublic.Username = user.Name
			}

			if localUser.Username == "" {
				localUser.UserPublic.Username = strings.Split(user.Email, "@")[0]
			}

			localUser = RegisterUser(localUser.Username)
			localUser.LinkedAccounts = append(localUser.LinkedAccounts, user)
			localUser.Email = user.Email

			if localUser.Email != "" {
				localUser.UserPublic.Verified = true
			}

		} else if cursor.Err() == nil {
			cursor.Decode(&localUser)

		} else {
			writer.Header().Set("Location", "/")
			writer.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		AddSessionToUser(&localSession, localUser)
		database.UpdateUserByID(localUser.UserPublic.UserID, localUser)
		session, _ := Store.Get(request, "sh-session")
		session.Values["session"] = localSession
		_ = session.Save(request, writer)

		if localUser.FinishedSignup {
			writer.Header().Set("Location", "/game/")
			writer.WriteHeader(http.StatusTemporaryRedirect)

		} else {
			writer.Header().Set("Location", "/oauth-select-username")
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}
	})
}
