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
//	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()

func Authenticate(request *http.Request) *types.Session {
	session, _ := Store.Get(request, "session")
	localSession, ok := session.Values["session"].(types.Session)

	if !ok {
		return nil
	}

	fmt.Println(localSession)

	cursor := database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
		"token": localSession.Token,
		"userID": localSession.UserID,
	})

	fmt.Println(cursor.Err(), "???")

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

func Render(tmplName string) http.Handler {
	fmt.Println("./views/" + tmplName + ".tmpl")
	tmpl := template.Must(template.ParseFiles("./views/" + tmplName + ".tmpl")).Funcs(template.FuncMap{
		"marshal": Marshal,
	})

	return http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
		localSession := Authenticate(request)
		localUser := types.User{}
		// fmt.Println("LS**********:", localSession)

		if localSession != nil {
			localUser.UserID = localSession.UserID

			cursor := database.MongoDB.Collection("Users").FindOne(ctx, bson.M{
				"userID": localUser.UserID,
			})

			// fmt.Println(cursor.Err())

			if cursor.Err() == nil {
				cursor.Decode(&localUser)
			}
		}

		// fmt.Println("LUUUUUUUUUUUU", localUser)

		data := struct {
			Game                       interface{}
			ProdCacheBustToken         string
			Username                   string
			Home                       interface{}
			Changelog                  interface{}
			Rules                      interface{}
			Howtoplay                  interface{}
			Stats                      interface{}
			Wiki                       interface{}
			Discord                    interface{}
			Github                     interface{}
			Tou                        interface{}
			About                      interface{}
			PrimaryColor               template.CSS
			SecondaryColor             template.CSS
			TertiaryColor              template.CSS
			BackgroundColor            template.CSS
			SecondaryBackgroundColor   template.CSS
			TertiaryBackgroundColor    template.CSS
			TextColor                  template.CSS
			SecondaryTextColor         template.CSS
			TertiaryTextColor          template.CSS
			GameSettings               interface{}
			Verified                   interface{}
			StaffRole                  string
			HasNotDismissedSignupModal interface{}
			IsTournamentMod            interface{}
			Blacklist                  interface{}
		} {
			Game: false,
			ProdCacheBustToken: CacheToken,
			Username: localUser.Username,
			Home: false,
			Changelog: false,
			Rules: false,
			Howtoplay: false,
			Stats: false,
			Wiki: false,
			Discord: false,
			Github: false,
			Tou: false,
			About: false,
			PrimaryColor: template.CSS("hsl(225, 73%, 57%)"),
			SecondaryColor: template.CSS("hsl(225, 48%, 57%)"),
			TertiaryColor: template.CSS("hsl(265, 73%, 57%)"),
			BackgroundColor: template.CSS("hsl(0, 0%, 0%)"),
			SecondaryBackgroundColor: template.CSS("hsl(0, 0%, 7%)"),
			TertiaryBackgroundColor: template.CSS("hsl(0, 0%, 14%)"),
			TextColor: template.CSS("hsl(0, 0%, 100%)"),
			SecondaryTextColor: template.CSS("hsl(0, 0%, 93%)"),
			TertiaryTextColor: template.CSS("hsl(0, 0%, 86%)"),
			GameSettings: types.GameSettings{
				CustomWidth: "",
				FontFamily: "",
			},
			Verified: struct{}{},
			StaffRole: localUser.StaffRole,
			HasNotDismissedSignupModal: struct{}{},
			IsTournamentMod: struct{}{},
			Blacklist: struct{}{},
		}

		if tmplName == "game" {
			data.Game = true

		} else if tmplName == "page-home" {
			data.Home = true
		}

		data.ProdCacheBustToken = CacheToken
		tmpl.Execute(writer, data)
	})
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

	router.HandleFunc("/online-playercount", func (writer http.ResponseWriter, request *http.Request) {
		data, _ := database.RedisDB.Get(ctx, "player-count").Result()
		utils.JSONResponse(writer, struct {
			Count string `json:"count"`
		} { data }, 200)
	})

	router.Handle("/socket.io/", socket.SetupSocketRoutes(io, store))

	router.HandleFunc("/{provider}-login", func (writer http.ResponseWriter, request *http.Request) {
		_, err := gothic.CompleteUserAuth(writer, request)

		if err == nil {
			writer.Header().Set("Location", "/game/")
			writer.WriteHeader(http.StatusTemporaryRedirect)

		} else {
			gothic.BeginAuthHandler(writer, request)
		}
	})

	router.HandleFunc("/auth/{provider}/callback", func (writer http.ResponseWriter, request *http.Request) {
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

		localUser := types.User{}

		// fmt.Println(cursor.Err(), user, user.NickName, user.Name, "**")

		if cursor.Err() == mongo.ErrNoDocuments {
			localUser.Username = user.NickName

			if localUser.Username == "" {
				localUser.Username = user.Name
			}

			// fmt.Println(localUser, "*", localUser.Username, "*", localUser.NickName, localUser.Name, "**")
			localUser.Username = strings.ReplaceAll(localUser.Username, " ", "-")
			// fmt.Println(localUser, "*", localUser.Username, "*", localUser.NickName, localUser.Name, "**")

			cursor = database.MongoDB.Collection("Users").FindOne(ctx, types.User{
				Username: localUser.Username,
			})

			if cursor.Err() == mongo.ErrNoDocuments {
				userID := uuid.NewString()

				cursor = database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
					"userID": userID,
				})

				for cursor.Err() != mongo.ErrNoDocuments {
					userID = uuid.NewString()

					cursor = database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
						"userID": userID,
					})
				}

				localUser.UserID = userID
				database.MongoDB.Collection("Users").InsertOne(ctx, localUser)

			} else {
				cursor.Decode(&localUser)
			}

		} else if cursor.Err() == nil {
			cursor.Decode(&localUser)

		} else {
			writer.Header().Set("Location", "/")
			writer.WriteHeader(http.StatusTemporaryRedirect)
		}

		localSession.UserID = localUser.UserID
		localSession.Expires = time.Now().Add(7 * 24 * time.Hour)
		database.MongoDB.Collection("Sessions").InsertOne(ctx, localSession)
		localUser.Sessions = append(localUser.Sessions, localSession)
		localUser.LinkedAccounts = append(localUser.LinkedAccounts, user)

		database.MongoDB.Collection("Users").UpdateOne(ctx, bson.M{
			"userID": localUser.UserID,

		}, bson.M{
			"$set": &localUser,
		})

		session, err := Store.Get(request, "session")
		session.Values["session"] = localSession
		err = session.Save(request, writer)
		writer.Header().Set("Location", "/game/")
		writer.WriteHeader(http.StatusTemporaryRedirect)
	})
}
