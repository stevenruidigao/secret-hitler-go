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

	user := database.GetUserByID(session.UserID)

	if user == nil || !user.FinishedSignup {
		return nil
	}

	UserMapMutex.Lock()
	UserMap[user.UserPublic.ID] = user
	UserMapMutex.Unlock()
	// fmt.Println("*User", user)

	return user
}

func Connect(socket socketio.Conn) {
	socket.SetContext("")
	fmt.Println("Connecting", socket.ID())

	socket.Emit("version", struct {
		Current interface{} `json:"current"`
	}{Current: struct {
		Number string `json:"number"`
	}{Number: constants.CURRENT_VERSION_NUMBER}})

	// fmt.Println("Sent version:", socket.ID())
	EmotesListMutex.RLock()
	socket.Emit("emoteList", EmotesList)
	// jsonData, err := json.Marshal(EmotesList)
	// fmt.Println("Sent emote list:", err)
	// fmt.Println("more details", string(jsonData), EmotesList)
	EmotesListMutex.RUnlock()

	socket.Emit("generalChats", GeneralChats)
	// SendUserList(socket)
	user := GetUser(socket)
	// fmt.Println("User", user)

	if user != nil {
		data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
		playerCount, _ := strconv.Atoi(data)
		data = strconv.Itoa(playerCount + 1)
		// fmt.Println("*", data, playerCount)
		database.RedisDB.Set(ctx, "playerCount", data, 0)
		user.Connections++
		UserMapMutex.Lock()
		UserMap[user.UserPublic.ID] = user
		// fmt.Println("UserMap", UserMap)
		UserMapMutex.Unlock()
		// fmt.Println("GameSettings", user.GameSettings)
		socket.Emit("gameSettings", user.GameSettings, len(user.GameSettings.Blacklist))
		UpdateUserStatus(&user.UserPublic, nil, "")
	}

	// fmt.Println("Updated player count")

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

	// fmt.Println("Connected:", socket.ID())
}

func Disconnect(socket socketio.Conn, reason string) {
	user := GetUser(socket)

	if user != nil {
		// database.RedisDB.Del(ctx, utils.GetCookie(socket.RemoteHeader().Get("cookie"), "session"))
		// fmt.Println("disconnected", reason)
		data, _ := database.RedisDB.Get(ctx, "playerCount").Result()
		playerCount, _ := strconv.Atoi(data)
		data = strconv.Itoa(playerCount - 1)
		// fmt.Println("**", data, playerCount)
		database.RedisDB.Set(ctx, "playerCount", data, 0)

		if user.Status != nil && user.Status.Type == "playing" {
			GameMapMutex.RLock()
			game := GameMap[user.Status.GameID]
			GameMapMutex.RUnlock()

			playerNumber := game.GamePublic.PlayerMap[user.ID]

			if playerNumber > 0 {
				// game.GamePublic.PlayerCount--
				game.GamePublic.PublicPlayerStates[playerNumber-1].Connected = false
				IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID, "gameUpdate", game.GamePublic)
			}
			// UpdateUserStatus(user.UserPublic, nil, "")
		}

		user.Connections--

		UserMapMutex.Lock()

		if user.Connections == 0 {
			delete(UserMap, user.UserPublic.ID)

		} else {
			UserMap[user.UserPublic.ID] = user
		}

		UserMapMutex.Unlock()

		IO.BroadcastToRoom("/", "aem", "userList", GetUserList(true))
		IO.BroadcastToRoom("/", "users", "userList", GetUserList(false))
	}
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
		go Connect(socket)

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

		if user != nil {
			// a, _ := utils.MarshalJSON(user.GameSettings)
			// fmt.Println("gameSettings:", user, "\n*", user.GameSettings, "\n**", user.GameSettings.Private, "\n***", user.GameSettings.Blacklist, "\n****", a)
			socket.Emit("gameSettings", user.GameSettings)

		} else {
			// socket.Emit("gameSettings", struct{}{})
		}
	})

	IO.OnEvent("/", "addNewGame", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("AAaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\n\n\n\n\n\n\n")
		fmt.Println("addNewGame", data)

		go AddNewGame(socket, data)

		// socket.Emit("hi", "hi")
	})

	IO.OnEvent("/", "addNewGameChat", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("addNewGameChat", data)

		go AddNewGameChat(socket, data)
	})

	// IO.OnEvent("/", "addNewGeneralChat", func(socket socketio.Conn, data interface{}) {
	// 	fmt.Printf("addNewGeneralChat type: %T\n", data)
	// })

	IO.OnEvent("/", "addNewGeneralChat", func(socket socketio.Conn, data map[string]interface{}) {
		// fmt.Println("addNewGeneralChat", data)
		go AddNewGeneralChat(socket, data)
	})

	IO.OnEvent("/", "getGameInfo", func(socket socketio.Conn, id string) {
		go SendGameInfo(socket, id)
		fmt.Println("Updated Game")
	})

	IO.OnEvent("/", "updateSeatedUser", func(socket socketio.Conn, data map[string]interface{}) {
		go UpdateSeatedUser(socket, data)
	})

	IO.OnEvent("/", "presidentSelectedChancellor", func(socket socketio.Conn, data map[string]interface{}) {
		go func() {
			user := GetUser(socket)

			if user != nil {
				id, ok := data["uid"].(string)

				if !ok {
					return
				}

				chancellorIndex, ok := data["chancellorIndex"].(float64)

				if !ok {
					return
				}

				GameMapMutex.RLock()
				game := GameMap[id]
				GameMapMutex.RUnlock()
				SelectChancellor(&user.UserPublic, game, int(chancellorIndex), false)
			}
		}()
	})

	IO.OnEvent("/", "selectedVoting", func(socket socketio.Conn, data map[string]interface{}) {
		go func() {
			user := GetUser(socket)

			if user != nil {
				id, ok := data["uid"].(string)

				if !ok {
					return
				}

				vote, ok := data["vote"].(bool)

				if !ok {
					return
				}

				GameMapMutex.RLock()
				game := GameMap[id]
				GameMapMutex.RUnlock()

				for i := range game.SeatedPlayers {
					if game.SeatedPlayers[i].ID == user.UserPublic.ID {
						SelectVote(&game.SeatedPlayers[i], game, vote, false)
					}
				}
			}
		}()
	})

	IO.OnEvent("/", "selectedPresidentPolicy", func(socket socketio.Conn, data map[string]interface{}) {
		go func() {
			user := GetUser(socket)

			if user != nil {
				id, ok := data["uid"].(string)

				if !ok {
					return
				}

				selection, ok := data["selection"].(float64)

				if !ok {
					return
				}

				GameMapMutex.RLock()
				game := GameMap[id]
				GameMapMutex.RUnlock()

				SelectPresidentPolicy(&user.UserPublic, game, int(selection), false)
			}
		}()
	})

	IO.OnEvent("/", "selectedChancellorPolicy", func(socket socketio.Conn, data map[string]interface{}) {
		go func() {
			user := GetUser(socket)

			if user != nil {
				id, ok := data["uid"].(string)

				if !ok {
					return
				}

				fmt.Println("CSelection id", id, ok)

				selection, ok := data["selection"].(float64)

				if !ok {
					return
				}

				GameMapMutex.RLock()
				game := GameMap[id]
				GameMapMutex.RUnlock()

				fmt.Println("CSelection", selection, ok)

				if selection == 1 {
					SelectChancellorPolicy(&user.UserPublic, game, 0, false)

				} else if selection == 3 {
					SelectChancellorPolicy(&user.UserPublic, game, 1, false)
				}
			}
		}()
	})

	IO.OnEvent("/", "hasSeenNewPlayerModal", func(socket socketio.Conn, message string) {
		go func() {
			fmt.Println("hasSeenNewPlayerModal:", message)
			user := GetUser(socket)
			user.DismissedSignupModal = true

			database.MongoDB.Collection("Users").UpdateOne(ctx, bson.M{
				"user.UserPublic.ID": user.UserPublic.ID,
			}, bson.M{
				"$set": &user,
			})
		}()
	})

	IO.OnEvent("/", "updateGameSettings", func(socket socketio.Conn, data map[string]interface{}) {
		go func() {
			user := GetUser(socket)

			// fmt.Println("Update Game Settings:", data["rightSidebarInGame"], user.GameSettings)

			if user != nil {
				for key, value := range data {
					switch key {
					case "enableRightSidebarInGame":
						user.GameSettings.RightSidebarInGame = value.(bool)

					case "enableTimestamps":
						user.GameSettings.Timestamps = value.(bool)
					}
				}

				database.UpdateUserByID(user.UserPublic.ID, user)
			}

			socket.Emit("gameSettings", user.GameSettings)
		}()
	})

	IO.OnError("/", func(socket socketio.Conn, err error) {
		fmt.Println("********************Error:", err)
	})

	IO.OnDisconnect("/", func(socket socketio.Conn, reason string) {
		go Disconnect(socket, reason)
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
