package routes

import (
	"secrethitler.io/database"
	"secrethitler.io/socket"
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"context"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
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
	router.HandleFunc("/howtoplay", func(writer http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.ParseFiles("./views/layout.tmpl", "./views/page-howtoplay.tmpl"))
		tmpl.Execute(writer, types.RenderData{})
	}).Methods("GET")
	router.HandleFunc("/account", func(writer http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.ParseFiles("./views/layout.tmpl", "./views/page-account.tmpl"))
		tmpl.Execute(writer, types.RenderData{Username: "example", Verified: true, Email: "example@example.com", DiscordUsername: "example", DiscordDiscriminator: "0000", GithubUsername: "example"})
	}).Methods("GET")

	router.Handle("/oauth-select-username", Render("page-new-username")).Methods("GET")

	router.HandleFunc("/oauth-select-username", func(writer http.ResponseWriter, request *http.Request) {
		session := Authenticate(request)

		if session == nil {
			return
		}

		user := database.GetUserByID(session.UserID)

		if user == nil {
			return
		}

		username := struct {
			Username string `bson:"username" json:"username"`
		}{}

		json.NewDecoder(request.Body).Decode(&username)

		UpdateUserUsername(user, username.Username)
	}).Methods("POST")

	router.HandleFunc("/account/signup", func(writer http.ResponseWriter, request *http.Request) {
		// {username: "a", password: "", password2: "", email: "", isPrivate: false, bypass: ""}

		signupOptions := struct {
			Username             string `bson:"username" json:"username"`
			Password             string `bson:"password" json:"password"`
			PasswordConfirmation string `bson:"passwordConfirmation" json:"password2"`
			Email                string `bson:"email" json:"email"`
			Private              bool   `bson:"private" json:"isPrivate"`
			BypassKey            string `bson:"bypassKey" json:"bypass"`
		}{}

		json.NewDecoder(request.Body).Decode(&signupOptions)

		cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.username": signupOptions.Username,
		})

		if cursor.Err() != mongo.ErrNoDocuments {
			utils.JSONResponse(writer, struct {
				Message string `bson:"message" json:"message"`
			}{
				Message: "That account already exists.",
			}, 401)

			return
		}

		bytes := make([]byte, 32)
		rand.Read(bytes)
		salt := hex.EncodeToString(bytes)

		user := &types.UserPrivate{
			UserPublic: types.UserPublic{
				Username: signupOptions.Username,
			},
			PasswordHash: utils.Argon2(signupOptions.Password, salt),
			Salt:         salt,
			Email:        signupOptions.Email,
			Local:        true,
		}

		user = RegisterUser(user)

		if user == nil {
			utils.JSONResponse(writer, struct {
				Message string `bson:"message" json:"message"`
			}{
				Message: "Something went wrong.",
			}, 500)
		}

		localSession := types.Session{
			Token: uuid.NewString(),
		}

		AddSessionToUser(&localSession, user)
		database.UpdateUserByID(user.UserPublic.ID, user)
		session, _ := Store.Get(request, "sh-session")
		session.Values["session"] = localSession
		session.Save(request, writer)
	})

	router.HandleFunc("/account/signin", func(writer http.ResponseWriter, request *http.Request) {
		loginOptions := struct {
			Username string `bson:"username" json:"username"`
			Password string `bson:"password" json:"password"`
		}{}

		json.NewDecoder(request.Body).Decode(&loginOptions)

		cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"userPublic.username": loginOptions.Username,
		})

		if cursor.Err() == mongo.ErrNoDocuments {
			utils.JSONResponse(writer, struct {
				Message string `bson:"message" json:"message"`
			}{
				Message: "That account doesn't exist.",
			}, 401)

			return
		}

		user := types.UserPrivate{
			UserPublic: types.UserPublic{
				Created:    time.Now(),
				EloOverall: 1600,
				EloSeason:  1600,
				Status:     nil,
				Profile: types.Profile{
					Created:       time.Now(),
					LastConnected: time.Now(),
					Badges:        []types.Badge{},
					RecentGames:   []types.RecentGame{},
				},
			},
			GameSettings: types.GameSettings{
				Blacklist: []string{},
			},
		}

		cursor.Decode(&user)

		if utils.Argon2(loginOptions.Password, user.Salt) != user.PasswordHash {
			utils.JSONResponse(writer, struct {
				Message string `bson:"message" json:"message"`
			}{
				Message: "That password is incorrect.",
			}, 401)

			return
		}

		localSession := types.Session{
			Token:  uuid.NewString(),
			UserID: user.UserPublic.ID,
		}

		AddSessionToUser(&localSession, &user)
		database.UpdateUserByID(user.UserPublic.ID, &user)
		session, _ := Store.Get(request, "sh-session")
		session.Values["session"] = localSession
		session.Save(request, writer)
	}).Methods("POST")

	router.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		session, _ := Store.Get(request, "sh-session")
		session.Values["session"] = nil
		session.Save(request, writer)
		writer.Header().Set("Location", "/")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	})

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

		var user types.UserPrivate
		cursor.Decode(&user)
		// str, _ := utils.MarshalJSON(types.Profile{})
		// fmt.Println("*&&&&&", user, user.Profile, ";;;;;;;;;;", user.Profile.RecentGames, len(user.Profile.RecentGames))
		utils.JSONResponse(writer, user.UserPublic.Profile, 200)
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

		cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
			"linkedAccounts.provider": user.Provider,
			"linkedAccounts.userid":   user.UserID,
		})

		localSession := Authenticate(request)

		localUser := &types.UserPrivate{
			UserPublic: types.UserPublic{
				Created:    time.Now(),
				EloOverall: 1600,
				EloSeason:  1600,
				Status:     nil,
				Profile: types.Profile{
					Created:       time.Now(),
					LastConnected: time.Now(),
					Badges:        []types.Badge{},
					RecentGames:   []types.RecentGame{},
				},
			},
			GameSettings: types.GameSettings{
				Blacklist: []string{},
			},
		}

		// fmt.Println(cursor.Err(), user, user.NickName, user.Name, user.Email, "**")

		if cursor.Err() == mongo.ErrNoDocuments {
			if localSession != nil {
				localUser := database.GetUserByID(localSession.UserID)

				if localUser != nil {
					localUser.LinkedAccounts = append(localUser.LinkedAccounts, user)

					if user.Email != "" && localUser.Email == "" {
						localUser.Email = user.Email
						localUser.Verified = true
					}

					database.UpdateUserByID(localUser.UserPublic.ID, localUser)
					writer.Header().Set("Location", "/game/")
					writer.WriteHeader(http.StatusTemporaryRedirect)

					return
				}
			}

			localUser.UserPublic.Username = user.NickName

			if localUser.Username == "" {
				localUser.UserPublic.Username = user.Name
			}

			if localUser.Username == "" {
				localUser.UserPublic.Username = strings.Split(user.Email, "@")[0]
			}

			localUser = RegisterUser(localUser)
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

		localSession = &types.Session{
			Token: uuid.NewString(),
		}

		cursor = database.MongoDB.Collection("Sessions").FindOne(ctx, localSession)

		for cursor.Err() != mongo.ErrNoDocuments {
			localSession.Token = uuid.NewString()
			cursor = database.MongoDB.Collection("Sessions").FindOne(ctx, localSession)
		}

		AddSessionToUser(localSession, localUser)
		database.UpdateUserByID(localUser.UserPublic.ID, localUser)
		session, _ := Store.Get(request, "sh-session")
		session.Values["session"] = localSession
		session.Save(request, writer)

		if localUser.FinishedSignup {
			writer.Header().Set("Location", "/game/")
			writer.WriteHeader(http.StatusTemporaryRedirect)

		} else {
			writer.Header().Set("Location", "/oauth-select-username")
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}
	})
}
