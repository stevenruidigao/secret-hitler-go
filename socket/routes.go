package socket

import (
	"secrethitler.io/constants"
	"secrethitler.io/database"
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"context"
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	//	"time"

	//	"github.com/go-redis/redis/v8"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/sessions"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()
var IO *socketio.Server
var tokens = map[string]string{}
var tokensMutex = sync.RWMutex{}

func Authenticate(socket socketio.Conn) *types.Session {
	tokensMutex.RLock()
	token := tokens[utils.GetCookie(socket.RemoteHeader().Get("Cookie"), "sh-session")]
	tokensMutex.RUnlock()

	cur := database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
		"token": token,
	})

	// fmt.Println("T", token, tokens, utils.GetCookie(socket.RemoteHeader().Get("Cookie"), "sh-session"))

	if cur.Err() == mongo.ErrNoDocuments {
		return nil
	}

	session := types.Session{}
	cur.Decode(&session)

	return &session
}

func GetUser(socket socketio.Conn) *types.UserPrivate {
	session := Authenticate(socket)
	// fmt.Println("Session", session)

	if session == nil {
		return nil
	}

	return database.GetUserByID(session.UserID)
}

func SetupSocketRoutes(io *socketio.Server, store *sessions.CookieStore) http.HandlerFunc {
	IO = io

	files, err := ioutil.ReadDir("./public/images/emotes")

	if err == nil {
		EmotesListMutex.Lock()

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".png") {
				EmotesList[":"+strings.TrimSuffix(file.Name(), ".png")+":"] = "/images/emotes/" + file.Name() + "?v=" + constants.CURRENT_VERSION_NUMBER
			}
		}

		EmotesListMutex.Unlock()
	}

	IO.OnConnect("/", func(socket socketio.Conn) error {
		socket.SetContext("")
		fmt.Println("Connecting", socket.ID())

		socket.Emit("version", struct {
			Current interface{} `json:"current"`
		}{Current: struct {
			Number string `json:"number"`
		}{Number: constants.CURRENT_VERSION_NUMBER}})

		fmt.Println("Sent version:", socket.ID())
		EmotesListMutex.RLock()
		socket.Emit("emoteList", EmotesList)
		// jsonData, err := json.Marshal(EmotesList)
		// fmt.Println("Sent emote list:", err)
		// fmt.Println("more details", string(jsonData), EmotesList)
		EmotesListMutex.RUnlock()
		// SendUserList(socket)
		user := GetUser(socket)
		// fmt.Println("User", user)

		if user != nil {
			data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
			playerCount, _ := strconv.Atoi(data)
			data = strconv.Itoa(playerCount + 1)
			fmt.Println("*", data, playerCount)
			database.RedisDB.Set(ctx, "playerCount", data, 0)
			UserMapMutex.Lock()
			UserMap[user.UserPublic.UserID] = user.UserPublic
			// fmt.Println("UserMap", UserMap)
			UserMapMutex.Unlock()
			socket.Emit("gameSettings", user.GameSettings)
			UpdateUserStatus(user.UserPublic, nil, "")
		}

		fmt.Println("Updated player count")

		if user != nil && user.StaffRole != "" && user.StaffRole != "altmod" && user.StaffRole != "veteran" {
			IO.JoinRoom("/", "aem", socket)

		} else {
			IO.JoinRoom("/", "users", socket)
		}

		IO.BroadcastToRoom("/", "aem", "userList", GetUserList(true))
		IO.BroadcastToRoom("/", "users", "userList", GetUserList(false))
		SendGameList(socket)

		// crash here
		// if user != nil {
		// 	socket.Emit("gameSettings", user.GameSettings)
		// }

		fmt.Println("Connected:", socket.ID())

		return nil
	})

	IO.OnEvent("/", "notice", func(socket socketio.Conn, message string) {
		fmt.Println("notice:", message)
		socket.Emit("reply", "have "+message)
	})

	IO.OnEvent("/", "connection", func(socket socketio.Conn, message string) {
		fmt.Println("connection:", message)
		socket.Emit("reply", "have "+message)
	})

	IO.OnEvent("/", "getUserGameSettings", func(socket socketio.Conn, message string) {
		user := GetUser(socket)

		if user == nil {
			socket.Emit("gameSettings", struct{}{})
			return
		}

		// a, _ := utils.MarshalJSON(user.GameSettings)
		// fmt.Println("gameSettings*", user, "*", user.GameSettings, "**", user.GameSettings.Private, "\n", a)
		socket.Emit("gameSettings", user.GameSettings)
	})

	IO.OnEvent("/", "addNewGame", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("AAaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\n\n\n\n\n\n\n")
		// fmt.Println("addNewGame", data)
		user := GetUser(socket)

		if user != nil {
			AddNewGame(socket, user.UserPublic, data)
		}

		socket.Emit("hi", "hi")
	})

	IO.OnEvent("/", "addNewGameChat", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("addNewGameChat", data)
		user := GetUser(socket)

		if user != nil {
			id, ok := data["uid"].(string)

			if ok {
				GameMapMutex.RLock()
				game := GameMap[id]
				GameMapMutex.RUnlock()

				AddNewGameChat(socket, user.UserPublic, data, game)
			}
		}
	})

	IO.OnEvent("/", "addNewGeneralChat", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("addNewGeneralChat", data)
		user := GetUser(socket)

		if user != nil {
			AddNewGeneralChat(socket, user.UserPublic, data)
		}
	})

	IO.OnEvent("/", "getGameInfo", func(socket socketio.Conn, id string) {
		user := GetUser(socket)

		fmt.Println("Game Update:", id)

		if user != nil {
			SendGameInfo(socket, &user.UserPublic, id)

		} else {
			SendGameInfo(socket, nil, id)
		}

		fmt.Println("Updated Game")
	})

	IO.OnEvent("/", "updateSeatedUser", func(socket socketio.Conn, data map[string]interface{}) {
		user := GetUser(socket)

		if user != nil {
			UpdateSeatedUser(socket, user.UserPublic, data)
		}
	})

	IO.OnEvent("/", "hasSeenNewPlayerModal", func(socket socketio.Conn, message string) {
		fmt.Println("hasSeenNewPlayerModal:", message)
		user := GetUser(socket)
		user.DismissedSignupModal = true

		database.MongoDB.Collection("Users").UpdateOne(ctx, bson.M{
			"user.UserPublicID": user.UserPublic.UserID,
		}, bson.M{
			"$set": &user,
		})
	})

	IO.OnEvent("/", "updateGameSettings", func(socket socketio.Conn, data map[string]interface{}) {
		user := GetUser(socket)

		// fmt.Println("Update Game Settings:", data["rightSidebarInGame"], user.GameSettings)

		if user != nil {
			for key, value := range data {
				switch key {
				case "enableRightSidebarInGame":
					user.GameSettings.RightSidebarInGame = value.(bool)
					break
				case "enableTimestamps":
					user.GameSettings.Timestamps = value.(bool)
					break
				}
			}

			database.UpdateUserByID(user.UserPublic.UserID, user)
		}
	})

	IO.OnError("/", func(socket socketio.Conn, err error) {
		fmt.Println("********************Error:", err)
	})

	IO.OnDisconnect("/", func(socket socketio.Conn, reason string) {
		user := GetUser(socket)

		if user != nil {
			// database.RedisDB.Del(ctx, utils.GetCookie(socket.RemoteHeader().Get("cookie"), "session"))
			fmt.Println("disconnected", reason)
			data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
			playerCount, _ := strconv.Atoi(data)
			data = strconv.Itoa(playerCount - 1)
			fmt.Println("**", data, playerCount)
			database.RedisDB.Set(ctx, "playerCount", data, 0)

			if user.Status.Type == "playing" {
				GameMapMutex.RLock()
				game := GameMap[user.Status.GameID]
				GameMapMutex.RUnlock()

				playerNumber, ok := game.GamePublic.GeneralGameSettings.Map[user.UserID].(int)

				if ok {
					game.GamePublic.PlayerCount--
					game.GamePublic.PublicPlayersState[playerNumber].Connected = false
					IO.BroadcastToRoom("/", "game-"+game.GeneralGameSettings.ID, "gameUpdate", game.GamePublic)
				}
				// UpdateUserStatus(user.UserPublic, nil, "")
			}

			UserMapMutex.Lock()
			delete(UserMap, user.UserPublic.UserID)
			UserMapMutex.Unlock()

			IO.BroadcastToRoom("/", "aem", "userList", GetUserList(true))
			IO.BroadcastToRoom("/", "users", "userList", GetUserList(false))
		}
	})

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "sh-session")
		localSession, ok := session.Values["session"].(types.Session)

		if ok {
			sessionCookie, _ := request.Cookie("sh-session")
			tokensMutex.Lock()
			tokens[sessionCookie.Value] = localSession.Token
			tokensMutex.Unlock()

			// fmt.Println("Session:", localSession.Token)
		}

		IO.ServeHTTP(writer, request)
	})
}
