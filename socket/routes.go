package socket

import (
	"secrethitler.io/database"
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"context"
	"fmt"
	"net/http"
	"strconv"
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
	token := tokens[utils.GetCookie(socket.RemoteHeader().Get("Cookie"), "session")]
	tokensMutex.RUnlock()

	cur := database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
		"token": token,
	})

	// fmt.Println("T", token)

	if cur.Err() == mongo.ErrNoDocuments {
		return nil
	}

	session := types.Session{}
	cur.Decode(&session)

	return &session
}

func GetUser(socket socketio.Conn) *types.UserPrivate {
	session := Authenticate(socket)

	if session == nil {
		return nil
	}

	return database.GetUserByID(session.UserID)
}

func SetupSocketRoutes(io *socketio.Server, store *sessions.CookieStore) http.HandlerFunc {
	IO = io

	IO.OnConnect("/", func(socket socketio.Conn) error {
		socket.SetContext("")
		fmt.Println("Connecting", socket.ID())

		socket.Emit("version", struct {
			Current interface{} `json:"current"`
		}{Current: struct {
			Number string `json:"number"`
		}{Number: "1.8.2"}})

		fmt.Println("Sent version:", socket.ID())

		// SendUserList(socket)
		user := GetUser(socket)

		if user != nil {
			data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
			playerCount, _ := strconv.Atoi(data)
			data = strconv.Itoa(playerCount + 1)
			fmt.Println("*", data, playerCount)
			database.RedisDB.Set(ctx, "playerCount", data, 0)
			UserListMutex.Lock()
			UserList = append(UserList, user.User)
			UserListMutex.Unlock()
			socket.Emit("gameSettings", user.GameSettings)
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

		a, _ := utils.MarshalJSON(user.GameSettings)
		fmt.Println("gameSettings*", user, "*", user.GameSettings, "**", user.GameSettings.Private, "\n", a)
		socket.Emit("gameSettings", user.GameSettings)
		/*socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)*/
	})

	IO.OnEvent("/", "addNewGame", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("AAaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\n\n\n\n\n\n\n")
		fmt.Println("addNewGame", data)
		user := GetUser(socket)

		if user != nil {
			AddNewGame(socket, user.User, data)
		}
		socket.Emit("hi", "hi")
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
			var i int
			UserListMutex.RLock()

			for i := range UserList {
				if UserList[i].UserID == user.UserID {
					break
				}
			}

			UserListMutex.RUnlock()
			UserListMutex.Lock()

			if i == len(UserList)-1 {
				UserList = append([]types.User{}, UserList[:i]...)

			} else {
				// error occurs here
				// UserList = append(append([]types.User{}, UserList[:i]...), UserList[i+1:]...)
			}

			UserListMutex.Unlock()

			IO.BroadcastToRoom("/", "aem", "userList", GetUserList(true))
			IO.BroadcastToRoom("/", "users", "userList", GetUserList(false))
		}
	})

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "session")
		localSession, ok := session.Values["session"].(types.Session)

		if ok {
			sessionCookie, _ := request.Cookie("session")
			tokensMutex.Lock()
			tokens[sessionCookie.Value] = localSession.Token
			tokensMutex.Unlock()
		}

		IO.ServeHTTP(writer, request)
	})
}
