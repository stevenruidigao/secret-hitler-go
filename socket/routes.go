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
var tokens = map[string]string{}
var tokensMutex = sync.RWMutex{}

func Authenticate(socket socketio.Conn) *types.Session {
	tokensMutex.RLock()
	token := tokens[utils.GetCookie(socket.RemoteHeader().Get("Cookie"), "session")]
	tokensMutex.RUnlock()

	cur := database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
		"token": token,
	})

	fmt.Println("T", token)

	if cur.Err() == mongo.ErrNoDocuments {
		return nil
	}

	session := types.Session{}
	cur.Decode(&session)

	return &session
}

func GetUser(socket socketio.Conn) *types.User {
	session := Authenticate(socket)

	if session == nil {
		return nil
	}

	return database.GetUserByID(session.UserID)
}

func SetupSocketRoutes(io *socketio.Server, store *sessions.CookieStore) http.HandlerFunc {
	io.OnConnect("/", func (socket socketio.Conn) error {
		socket.SetContext("")

		socket.Emit("version", struct {
			Current interface{} `json:"current"`

		} { Current: struct { Number string `json:"number"` } { Number: "1.8.2" } })

		SendUserList(socket)
		user := GetUser(socket)

		if user != nil {
			data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
			playerCount, _ := strconv.Atoi(data)
			data = strconv.Itoa(playerCount + 1)
			// fmt.Println("*", data, playerCount)
			database.RedisDB.Set(ctx, "playerCount", data, 0)
			UserList = append(UserList, *user)
		}

		return nil
	})

	io.OnEvent("/", "notice", func (socket socketio.Conn, message string) {
		fmt.Println("notice:", message)
		socket.Emit("reply", "have "+message)
	})

	io.OnEvent("/", "connection", func (socket socketio.Conn, message string) {
		fmt.Println("connection:", message)
		socket.Emit("reply", "have "+message)
	})

	io.OnEvent("/", "getUserGameSettings", func (socket socketio.Conn, message string) {
		user := GetUser(socket)

		if user == nil {
			socket.Emit("gameSettings", struct {} {})
			return
		}

		fmt.Println("gameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings)
	})

	io.OnError("/", func (socket socketio.Conn, err error) {
		fmt.Println("meet error:", err)
	})

	io.OnDisconnect("/", func (socket socketio.Conn, reason string) {
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

			for i = 0; i < len(UserList); i ++ {
				if UserList[i].UserID == user.UserID {
					break;
				}
			}

			UserList = append(append([]types.User{}, UserList[:i]...), UserList[i+1:]...)
		}
	})

	return http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "session")
		localSession, ok := session.Values["session"].(types.Session)

		if ok {
			sessionCookie, _ := request.Cookie("session")
			tokensMutex.Lock()
			tokens[sessionCookie.Value] = localSession.Token
			tokensMutex.Unlock()
		}

		io.ServeHTTP(writer, request)
	})
}
