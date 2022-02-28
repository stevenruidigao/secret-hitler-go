package socket

import (
	"secrethitler.io/constants"
	"secrethitler.io/database"
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"math"
	// "strconv"
	"sync"
	"time"

	"github.com/googollee/go-socket.io"

	"go.mongodb.org/mongo-driver/bson"
)

/*addNewGame map[blindMode:false disableGamechat:false disableObserver:false disableObserverLobby:false eloSliderValue:<nil> excludedPlayerCount:[5 6 8 9 10] experiencedMode:true flag:none flappyMode:false flappyOnlyMode:false gameName:New Game gameType:ranked isTourny:false isVerifiedOnly:true maxPlayersCount:7 minPlayersCount:7 playerChats:enabled privateAnonymousRemakes:false privatePassword:false rainbowgame:false rebalance6p:false rebalance7p:false rebalance9p2f:false timedMode:false unlistedGame:false]*/

func AddNewGame(socket socketio.Conn, user *types.UserPrivate, data map[string]interface{}) {
	// gameSettings := data["gameSettings"]
	if user == nil {
		return
	}

	currentTime := time.Now()

	fmt.Println("userstatus", (user.UserPublic.Status != nil && user.UserPublic.Status.Type != "none"))

	if time.Since(user.TimeLastGameCreated) < time.Second*10 || (user.UserPublic.Status != nil && user.UserPublic.Status.Type != "none") {
		// fmt.Println("^^^", time.Now(), user.TimeLastGameCreated, time.Since(user.TimeLastGameCreated), "*", user.Status.Type, "*", user.Status.Type != "none", "*")
		return
	}

	fmt.Println("Creating new game")
	fmt.Println("Constants", constants.ALPHANUMERIC)
	// fmt.Printf("%T %f\n", data["customGameSettings"], data["customGameSettings"])

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		var ok bool
	// 		err, ok := r.(error)
	// 		if !ok {
	// 			fmt.Println("*****************************************pkg: %v", r, err)
	// 		}
	// 	}
	// }()

	dataExcludedPlayerCount := data["excludedPlayerCount"].([]interface{})
	excludedPlayerCount := []int{}

	for _, playerCount := range dataExcludedPlayerCount {
		excludedPlayerCount = append(excludedPlayerCount, int(playerCount.(float64)))
	}

	rainbowGame, ok := data["rainbowgame"].(bool)

	if !ok {
		rainbowGame = false
	}

	timer, ok := data["timedMode"].(float64)

	if !ok {
		timer = -1
	}

	fmt.Println("timer", timer)

	casualGame, ok := data["casualGame"].(bool)

	if !ok {
		casualGame = false
	}

	fmt.Println("casualGame", casualGame)

	practiceGame, ok := data["practiceGame"].(bool)

	if !ok {
		practiceGame = false
	}

	fmt.Println("practiceGame", practiceGame)

	remade, ok := data["isRemade"].(bool)

	if !ok {
		remade = false
	}

	fmt.Println("remade", remade)

	eloMinimum, ok := data["eloSliderValue"].(float64)

	if !ok {
		eloMinimum = -1
	}

	fmt.Println("eloMinimum", eloMinimum)

	// customGameSettings := map[string]interface{}{}
	customGameSettings, ok := data["customGameSettings"].(types.CustomGameSettings)

	if !ok {
		customGameSettings = types.CustomGameSettings{}
	}

	fmt.Println("customGameSettings", customGameSettings)

	// fmt.Println("excludedPlayerCount", excludedPlayerCount)
	// fmt.Println("experiencedMode", data["experiencedMode"].(bool))
	// fmt.Println(data["privatePassword"] != nil, data["customGameSettings"], data["customGameSettings"].(map[string]interface{}), data["casualGame"])

	gamePublic := types.GamePublic{
		Date: currentTime,
		GameState: types.GameState{
			PreviousElectedGovernment: []int{-1, -1},
			UndrawnPolicyCount:        17,
			DiscardedPolicyCount:      0,
			PresidentIndex:            -1,
		},
		Chats:   []types.PlayerChat{},
		Guesses: map[string]string{},
		GeneralGameSettings: types.GeneralGameSettings{
			WhitelistedPlayers:      []string{},
			ID:                      utils.GenerateCombination(4, "", true),
			Name:                    data["gameName"].(string),
			Flag:                    "none",
			MinPlayersCount:         int(data["minPlayersCount"].(float64)),
			ExcludedPlayerCount:     excludedPlayerCount,
			MaxPlayersCount:         int(data["maxPlayersCount"].(float64)),
			Status:                  "Waiting for more players...", //
			Experienced:             data["experiencedMode"].(bool),
			PlayerChats:             "enabled", //data["playerChats"].(string),
			VerifiedOnly:            data["isVerifiedOnly"].(bool),
			DisableObserverLobby:    data["disableObserverLobby"].(bool),
			DisableObserver:         data["disableObserver"].(bool),
			Tourny:                  false,
			LastModPing:             0,
			ChatReplTime:            []int{},
			DisableGamechat:         data["disableGamechat"].(bool),
			Rainbow:                 rainbowGame,
			Blind:                   data["blindMode"].(bool),
			Timer:                   int(timer),
			Flappy:                  data["flappyMode"].(bool),
			FlappyOnly:              data["flappyOnlyMode"].(bool),
			Casual:                  casualGame,
			Practice:                practiceGame,
			Rebalance6p:             data["rebalance6p"].(bool),
			Rebalance7p:             data["rebalance7p"].(bool),
			Rebalance9p2f:           data["rebalance9p2f"].(bool),
			Unlisted:                data["unlistedGame"].(bool),
			Private:                 false || !data["unlistedGame"].(bool) && data["privatePassword"] != false, //*/
			PrivateAnonymousRemakes: data["privateAnonymousRemakes"].(bool),
			PrivateOnly:             false,
			ElectionCount:           0,
			Remade:                  remade,
			EloMinimum:              int(eloMinimum),
			TimeCreated:             currentTime,
			Usernames:               []string{user.UserPublic.Username},
			CustomCardback:          []string{},
			CustomCardbackUID:       []string{},
			// Players: []types.Player{
			// 	types.Player{
			// 		UserPublic: user.UserPublic,
			// 		Connected:  true,
			// 	},
			// },
			SeatedCount: 1,
			// Map:             map[string]interface{}{user.UserPublic.ID: 0},
			// Mutex:           &sync.RWMutex{},
			GameCreatorName: user.UserPublic.Username,
			GameStatus:      "notStarted",
		},
		CustomGameSettings: customGameSettings,
		/*type PlayerState struct {
			UserID		   string    `bson:"userID" json:"userID"`
			Connected                bool   `bson:"connected" json:"connected"`
			LeftGame                 bool   `bson:"leftGame"  json:"leftGame"`
			CustomCardback           string `bson:"customCardback" json:"customCardback"`
			CustomCardbackID         string `bson:"customCardbackID" json:"customCardbackUid"`
			IsLoader                 bool   `bson:"isLoader" json:"isLoader"`
			IsRemakeVoting           bool   `bson:"isRemakeVoting" json:"isRemakeVoting"`
			PingTime                 int64  `bson:"pingTime" json:"pingTime"`
			Username                 string `bson:"username" json:"userName"`
			PreviousGovernmentStatus string `bson:"previousGovernmentStatus" json:"previousGovernmentStatus"`
			GovernmentStatus         string `bson:"governmentStatus" json:"governmentStatus"`
		}*/
		PublicPlayerStates: []types.PlayerState{
			types.PlayerState{
				UserPublic: user.UserPublic,
				Connected:  true,
				LeftGame:   false,
				// CustomCardback:           "",
				// CustomCardbackID:         "",
				// Loader:                   false,
				// RemakeVoting:             false,
				// PingTime:                 0,
				// PreviousGovernmentStatus: "",
				// GovernmentStatus:         "",
			},
		},
		// PlayerStates:     []PlayerState{},
		CardFlingerState: []interface{}{},
		TrackState: types.TrackState{
			LiberalPolicyCount:   0,
			FascistPolicyCount:   0,
			ElectionTrackerCount: 0,
			EnactedPolicies:      []types.Policy{},
		},
		PlayerCount:             0,
		PlayerMap:               map[string]int{user.UserPublic.ID: 0},
		ChatMutex:               &sync.RWMutex{},
		PublicPlayerStatesMutex: &sync.RWMutex{},
	}

	fmt.Println("Game object created")

	playerCounts := []int{}

	for playerCount := int(math.Round(math.Max(float64(gamePublic.GeneralGameSettings.MinPlayersCount), 5))); playerCount <= int(math.Round(math.Min(float64(gamePublic.GeneralGameSettings.MaxPlayersCount), 10))); playerCount++ {
		var element int
		for _, element = range excludedPlayerCount {
			if element == playerCount {
				break
			}
		}

		if element != playerCount {
			playerCounts = append(playerCounts, playerCount)
		}
	}

	fmt.Println("playerCounts", playerCounts)

	if len(playerCounts) == 0 {
		return
	}

	gamePublic.PlayerCounts = playerCounts
	gamePublic.PlayerCount = len(gamePublic.PublicPlayerStates)

	if gamePublic.GeneralGameSettings.Private {
		fmt.Println("private", gamePublic.GeneralGameSettings.Private, data["unlistedGame"], data["privatePassword"], data["privatePassword"] != false)
		gamePublic.GeneralGameSettings.Name = "Private Game"
	}

	if gamePublic.GeneralGameSettings.Name == "" {
		gamePublic.GeneralGameSettings.Name = "New Game"
	}

	fmt.Println("Name length", len(gamePublic.GeneralGameSettings.Name), gamePublic.GeneralGameSettings.Name)

	if len(gamePublic.GeneralGameSettings.Name) > 20 {
		fmt.Println("Name too long")
		gamePublic.GeneralGameSettings.Name = gamePublic.GeneralGameSettings.Name[:int(math.Round(math.Min(20, float64(len(gamePublic.GeneralGameSettings.Name)))))]
	}

	eloSliderValue, ok := data["eloSliderValue"].(float64)

	if !ok {
		eloSliderValue = -1
	}

	if data["eloSliderValue"] != nil && (user.UserPublic.EloSeason < eloSliderValue || user.UserPublic.EloOverall < eloSliderValue) {
		return
	}

	if gamePublic.CustomGameSettings.Enabled {

	} /* else {
		gamePublic.CustomGameSettings = map[string]interface{}{}
		gamePublic.CustomGameSettings.Enabled = false
	}*/

	if data["isTourny"] != false {
		gamePublic.GeneralGameSettings.ID += "Tourny"
	}

	gamePublic.ID = gamePublic.GeneralGameSettings.ID
	user.UserPublic.TimeLastGameCreated = currentTime

	database.MongoDB.Collection("Users").UpdateOne(ctx, bson.M{
		"userID": user.UserPublic.ID,
	}, bson.M{
		"$set": user,
	})

	fmt.Println("Update user")

	gamePrivate := types.GamePrivate{
		GamePublic:              gamePublic,
		Reports:                 struct{}{},
		UnseatedGameChats:       []types.PlayerChat{},
		CommandChats:            []types.PlayerChat{},
		ReplayGameChats:         []types.PlayerChat{},
		Lock:                    struct{}{},
		VotesPeeked:             false,
		RemakeVotesPeeked:       false,
		InvIndex:                -1,
		HiddenInfoChat:          []types.PlayerChat{},
		HiddenInfoSubscriptions: []interface{}{},
		HiddenInfoShouldNotify:  true,
		GameCreatorName:         user.UserPublic.Username,
		GameCreatorID:           user.UserPublic.ID,
		GameCreatorBlacklist:    []string{},
		SocketMap: map[string][]socketio.Conn{
			user.UserPublic.ID: []socketio.Conn{socket},
		},
		SocketMapMutex: &sync.RWMutex{},
	}

	gamePrivate.GamePublic.GeneralGameSettings.Status = DisplayWaitingForPlayers(&gamePrivate)
	privatePassword, ok := data["privatePassword"].(string)

	if ok {
		gamePrivate.PrivatePassword = privatePassword
		gamePrivate.GamePublic.GeneralGameSettings.Private = true
	}

	fmt.Println("Created private game object")

	GameMapMutex.Lock()
	GameMap[gamePublic.ID] = &gamePrivate
	GameMapMutex.Unlock()

	GameMapMutex.RLock()
	database.RedisDB.Set(ctx, "gamesMap", GameMap, 0)
	GameMapMutex.RUnlock()

	fmt.Println("Updated game map")

	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	IO.JoinRoom("/", "game-"+gamePrivate.GamePublic.GeneralGameSettings.ID, socket)
	// IO.JoinRoom("/", "game-"+gamePrivate.GamePublic.GeneralGameSettings.ID+"-observer", socket)
	IO.JoinRoom("/", "game-"+gamePrivate.GamePublic.GeneralGameSettings.ID+"-"+user.UserPublic.ID, socket)
	socket.Emit("updateSeatForUser")
	// a, _ := utils.MarshalJSON(gamePrivate.GamePublic)
	// fmt.Println("gamePublic", gamePrivate.GamePublic, a)
	socket.Emit("gameUpdate", gamePrivate.GamePublic)
	// fmt.Println("*status", gamePrivate.GamePublic.GeneralGameSettings.Status)
	socket.Emit("joinGameRedirect", gamePrivate.GamePublic.GeneralGameSettings.ID)
	fmt.Println("gamePrivate", gamePrivate)
}

func AddNewGameChat(socket socketio.Conn, user *types.UserPublic, data map[string]interface{}, game *types.GamePrivate) {
	if game == nil || game.GamePublic.ChatMutex == nil {
		return
	}

	fmt.Println("AddNewGameChat", data)
	chat, ok := data["chat"].(string)
	// mutex := game.GamePublic.GeneralGameSettings.Mutex
	// fmt.Println("*", mutex)

	if ok {
		game.GamePublic.ChatMutex.Lock()
		fmt.Println("still adding...")

		game.GamePublic.Chats = append(game.GamePublic.Chats, types.PlayerChat{
			Username:  user.Username,
			UserID:    user.ID,
			Chat:      chat, //[]types.GameChat{types.GameChat{Text: chat}},
			StaffRole: user.StaffRole,
			Timestamp: time.Now(),
			GameID:    game.ID,
		})

		game.GamePublic.ChatMutex.Unlock()
		GameMapMutex.Lock()
		GameMap[game.GamePublic.ID] = game
		GameMapMutex.Unlock()

		fmt.Println("Added new game chat")
		game.GamePublic.ChatMutex.RLock()
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.ID, "playerChatUpdate", game.GamePublic.Chats[len(game.GamePublic.Chats)-1])
		game.GamePublic.ChatMutex.RUnlock()
		// IO.BroadcastToRoom("/", "game-"+game.ID, "gameUpdate", game)
	}
}

func AddNewGeneralChat(socket socketio.Conn, user *types.UserPublic, data map[string]interface{}) {
	fmt.Println("AddNewGeneralChat", data)

	chat, ok := data["chat"].(string)

	if ok {
		GeneralChatsMutex.Lock()
		GeneralChats.List = append(GeneralChats.List, types.GeneralChat{
			Username:  user.Username,
			UserID:    user.ID,
			Message:   chat,
			StaffRole: user.StaffRole,
			Timestamp: time.Now(),
		})
		GeneralChatsMutex.Unlock()

		fmt.Println("Added new general chat")
		IO.BroadcastToRoom("/", "users", "generalChats", GeneralChats)
		IO.BroadcastToRoom("/", "aem", "generalChats", GeneralChats)
	}
}

func UpdateSeatedUser(socket socketio.Conn, user *types.UserPublic, data map[string]interface{}) {
	if user == nil {
		return
	}

	id, ok := data["uid"].(string)

	if !ok {
		return
	}

	game := GameMap[id]

	fmt.Println("Attempted Join", len(game.GamePublic.PublicPlayerStates), game.GamePublic.PlayerCounts[len(game.GamePublic.PlayerCounts)-1])

	game.GamePublic.PublicPlayerStatesMutex.RLock()
	fmt.Println("Player Check", len(game.GamePublic.PublicPlayerStates), game.GamePublic.PlayerCounts[len(game.GamePublic.PlayerCounts)-1])

	if len(game.GamePublic.PublicPlayerStates) == game.GamePublic.PlayerCounts[len(game.GamePublic.PlayerCounts)-1] {
		return
	}

	for i := range game.GamePublic.PublicPlayerStates {
		if game.GamePublic.PublicPlayerStates[i].UserPublic.ID == user.ID {
			// return
		}
	}

	game.GamePublic.PublicPlayerStatesMutex.RUnlock()

	if game.GameCreatorBlacklist != nil {
		for _, id := range game.GameCreatorBlacklist {
			if id == user.ID {
				socket.Emit("gameJoinStatusUpdate", bson.M{
					"status": "blacklisted",
				})

				return
			}
		}
	}

	// game.GamePublic.GeneralGameSettings.Players = append(game.GamePublic.GeneralGameSettings.Players, types.Player{
	// 	UserPublic: *user,
	// 	Connected:  true,
	// })

	game.GamePublic.PublicPlayerStatesMutex.Lock()

	game.GamePublic.PublicPlayerStates = append(game.GamePublic.PublicPlayerStates, types.PlayerState{
		UserPublic: *user,
		// Socket:     socket,
		// UserID:                   user.ID,
		Connected: true,
		LeftGame:  false,
		// CustomCardback:   "",
		// CustomCardbackID: "",
		// Loader:           false,
		// RemakeVoting:     false,
		// PingTime:         0,
		// Username:                 user.Username,
		// PreviousGovernmentStatus: "",
		// GovernmentStatus:         "",
	})

	game.GamePublic.GeneralGameSettings.Usernames = append(game.GamePublic.GeneralGameSettings.Usernames, user.Username)
	game.GamePublic.PublicPlayerStatesMutex.Unlock()

	// game.GamePublic.GeneralGameSettings.Map[user.ID] = game.GamePublic.PlayerCount
	game.GamePublic.PlayerMap[user.ID] = game.GamePublic.PlayerCount
	game.GamePublic.PlayerCount = len(game.GamePublic.PublicPlayerStates)
	game.GamePublic.GeneralGameSettings.Status = DisplayWaitingForPlayers(game)
	game.GamePublic.GeneralGameSettings.SeatedCount++

	game.SocketMapMutex.Lock()
	game.SocketMap[user.ID] = append(game.SocketMap[user.ID], socket)
	game.SocketMapMutex.Unlock()

	GameMapMutex.Lock()
	GameMap[game.GamePublic.ID] = game
	GameMapMutex.Unlock()

	GameMapMutex.RLock()
	database.RedisDB.Set(ctx, "gamesMap", GameMap, 0)
	GameMapMutex.RUnlock()

	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	IO.JoinRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID, socket)
	IO.LeaveRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-observer", socket)
	IO.JoinRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+user.ID, socket)
	socket.Emit("updateSeatForUser")
	IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID, "gameUpdate", game.GamePublic)
	UpdateUserStatus(user, &game.GamePublic, "playing")
	// socket.Emit("gameUpdate", game.GamePublic)
}
