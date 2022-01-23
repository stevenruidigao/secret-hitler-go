package socket

import (
	"secrethitler.io/database"
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"context"
	"fmt"
	"net/http"
	"strconv"
//	"time"

//	"github.com/go-redis/redis/v8"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/sessions"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()
var tokens = map[string]string{}

func Authenticate(socket socketio.Conn) *types.Session {
	token := tokens[utils.GetCookie(socket.RemoteHeader().Get("Cookie"), "id")]

	cur := database.MongoDB.Collection("Sessions").FindOne(ctx, bson.M{
		"token": token,
	})

	if cur.Err() == mongo.ErrNoDocuments {
		return nil
	}

	session := types.Session{}
	cur.Decode(&session)

	return &session
}

func SetupSocketRoutes(io *socketio.Server, store *sessions.CookieStore) http.HandlerFunc {
        io.OnConnect("/", func (socket socketio.Conn) error {
                socket.SetContext("")

                socket.Emit("version", struct {
                        Current interface{} `json:"current"`

                }{struct {
                        Number string `json:"number"`
                }{"1.8.2"}})

                data, _ := database.RedisDB.Get(ctx, "player-count").Result()
                playerCount, _ := strconv.Atoi(data)
                data = strconv.Itoa(playerCount + 1)
                fmt.Println("*", data, playerCount)
                database.RedisDB.Set(ctx, "player-count", data, 0)

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
                fmt.Println("getsettings:", message)
                socket.Emit("reply", "have "+message)
        })

        io.OnError("/", func (socket socketio.Conn, err error) {
                fmt.Println("meet error:", err)
        })

        io.OnDisconnect("/", func (socket socketio.Conn, reason string) {
		database.RedisDB.Del(ctx, utils.GetCookie(socket.RemoteHeader().Get("cookie"), "session"))
                fmt.Println("closed", reason)
                data, _ := database.RedisDB.Get(ctx, "player-count").Result()
                playerCount, _ := strconv.Atoi(data)
                data = strconv.Itoa(playerCount - 1)
                fmt.Println("**", data, playerCount)
                database.RedisDB.Set(ctx, "player-count", data, 0)
        })

	return http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "session")
		localSession, ok := session.Values["session"].(types.Session)

		if ok {
			sessionCookie, _ := request.Cookie("session")
			tokens[sessionCookie.Value] = localSession.Token
		}

		io.ServeHTTP(writer, request)
	})
}
