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

func AddNewGame(socket socketio.Conn, user types.UserPublic, data map[string]interface{}) {
	// gameSettings := data["gameSettings"]

	currentTime := time.Now()

	if time.Since(user.TimeLastGameCreated) < time.Second*10 || user.Status.Type != "none" {
		fmt.Println(time.Now(), user.TimeLastGameCreated, time.Since(user.TimeLastGameCreated), "*", user.Status.Type, "*", user.Status.Type != "none", "*")
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

	newGame := types.GamePublic{
		Date: currentTime,
		GameState: types.GameState{
			PreviousElectedGovernment: []int{},
			UndrawnPolicyCount:        17,
			DiscardedPolicyCount:      0,
			PresidentIndex:            -1,
		},
		Chats:   []types.GameChat{},
		Guesses: map[string]string{},
		GeneralGameSettings: types.GeneralGameSettings{
			WhitelistedPlayers:      []string{},
			ID:                      utils.GenerateCombination(4, "", true),
			Name:                    data["gameName"].(string),
			Flag:                    "none",
			MinPlayersCount:         int(data["minPlayersCount"].(float64)),
			ExcludedPlayerCount:     excludedPlayerCount,
			MaxPlayersCount:         int(data["maxPlayersCount"].(float64)),
			Status:                  "Waiting for more players...",
			Experienced:             data["experiencedMode"].(bool),
			PlayerChats:             "enabled",
			VerifiedOnly:            data["isVerifiedOnly"].(bool),
			DisableObserverLobby:    data["disableObserverLobby"].(bool),
			DisableObserver:         data["disableObserver"].(bool),
			Tourny:                  false,
			LastModPing:             0,
			ChatReplTime:            []int{},
			DisableGamechat:         data["disableGamechat"].(bool),
			Rainbow:                 data["rainbowgame"].(bool),
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
			Private:                 false || !data["unlistedGame"].(bool) && data["privatePassword"] != false,
			PrivateAnonymousRemakes: data["privateAnonymousRemakes"].(bool),
			PrivateOnly:             false,
			ElectionCount:           0,
			Remade:                  remade,
			EloMinimum:              int(eloMinimum),
			TimeCreated:             currentTime,
			Usernames:               []string{user.Username},
			CustomCardback:          []string{},
			CustomCardbackUID:       []string{},
			Players:                 []types.UserPublic{user},
			SeatedCount:             1,
			Map:                     map[string]interface{}{user.UserID: 0},
			Mutex:                   &sync.RWMutex{},
			GameCreatorName:         user.Username,
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
		PublicPlayersState: []types.PlayerState{
			types.PlayerState{
				UserID:                   user.UserID,
				Connected:                true,
				LeftGame:                 false,
				CustomCardback:           "",
				CustomCardbackID:         "",
				IsLoader:                 false,
				IsRemakeVoting:           false,
				PingTime:                 0,
				Username:                 user.Username,
				PreviousGovernmentStatus: "",
				GovernmentStatus:         "",
			},
		},
		PlayersState:     []interface{}{},
		CardFlingerState: []interface{}{},
		TrackState: types.TrackState{
			LiberalPolicyCount:   0,
			FascistPolicyCount:   0,
			ElectionTrackerCount: 0,
			EnactedPolicies:      []string{},
		},
		PlayerCount: 0,
	}

	playerCounts := []int{}

	for playerCount := int(math.Round(math.Max(float64(newGame.GeneralGameSettings.MinPlayersCount), 5))); playerCount <= int(math.Round(math.Min(float64(newGame.GeneralGameSettings.MaxPlayersCount), 10))); playerCount++ {
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

	newGame.PlayerCounts = playerCounts
	newGame.PlayerCount = len(newGame.PublicPlayersState)
	newGame.GeneralGameSettings.Status = DisplayWaitingForPlayers(&newGame)

	if newGame.GeneralGameSettings.Private {
		fmt.Println("private", newGame.GeneralGameSettings.Private, data["unlistedGame"], data["privatePassword"], data["privatePassword"] != false)
		newGame.GeneralGameSettings.Name = "Private Game"
	}

	if newGame.GeneralGameSettings.Name == "" {
		newGame.GeneralGameSettings.Name = "New Game"
	}

	fmt.Println("Name length", len(newGame.GeneralGameSettings.Name), newGame.GeneralGameSettings.Name)

	if len(newGame.GeneralGameSettings.Name) > 20 {
		fmt.Println("Name too long")
		newGame.GeneralGameSettings.Name = newGame.GeneralGameSettings.Name[:int(math.Round(math.Min(20, float64(len(newGame.GeneralGameSettings.Name)))))]
	}

	eloSliderValue, ok := data["eloSliderValue"].(float64)

	if !ok {
		eloSliderValue = -1
	}

	if data["eloSliderValue"] != nil && (user.EloSeason < eloSliderValue || user.EloOverall < eloSliderValue) {
		return
	}

	if newGame.CustomGameSettings.Enabled {

	} /* else {
		newGame.CustomGameSettings = map[string]interface{}{}
		newGame.CustomGameSettings.Enabled = false
	}*/

	if data["isTourny"] != false {
		newGame.GeneralGameSettings.ID += "Tourny"
	}

	newGame.ID = newGame.GeneralGameSettings.ID
	user.TimeLastGameCreated = currentTime

	database.MongoDB.Collection("Users").UpdateOne(ctx, bson.M{
		"userID": user.UserID,
	}, bson.M{
		"$set": user,
	})

	privateGame := types.GamePrivate{
		GamePublic:              newGame,
		Reports:                 struct{}{},
		UnseatedGameChats:       []types.GameChat{},
		CommandChats:            []types.GameChat{},
		ReplayGameChats:         []types.GameChat{},
		Lock:                    struct{}{},
		VotesPeeked:             false,
		RemakeVotesPeeked:       false,
		InvIndex:                -1,
		HiddenInfoChat:          []types.GameChat{},
		HiddenInfoSubscriptions: []interface{}{},
		HiddenInfoShouldNotify:  true,
		GameCreatorName:         user.Username,
		GameCreatorID:           user.UserID,
		GameCreatorBlacklist:    []string{},
	}

	privatePassword, ok := data["privatePassword"].(string)

	if ok {
		privateGame.PrivatePassword = privatePassword
		newGame.GeneralGameSettings.Private = true
	}

	GameMapMutex.Lock()
	GameMap[privateGame.GamePublic.ID] = privateGame
	GameMapMutex.Unlock()

	GameMapMutex.RLock()
	database.RedisDB.Set(ctx, "gamesMap", GameMap, 0)
	GameMapMutex.RUnlock()

	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	IO.JoinRoom("/", "game-"+newGame.GeneralGameSettings.ID, socket)
	socket.Emit("updateSeatForUser")
	socket.Emit("gameUpdate", newGame)
	socket.Emit("joinGameRedirect", newGame.GeneralGameSettings.ID)
	fmt.Println("newGame", newGame)
}

func AddNewGameChat(socket socketio.Conn, user *types.UserPublic, data map[string]interface{}, game *types.GamePrivate) {
	if game == nil || game.GamePublic.GeneralGameSettings.Mutex == nil {
		return
	}

	fmt.Println("AddNewGameChat", data)
	chat, ok := data["chat"].(string)
	// mutex := game.GamePublic.GeneralGameSettings.Mutex
	// fmt.Println("*", mutex)

	if ok {
		game.GamePublic.GeneralGameSettings.Mutex.Lock()
		fmt.Println("still adding...")

		game.GamePublic.Chats = append(game.GamePublic.Chats, types.GameChat{
			Username:  user.Username,
			UserID:    user.UserID,
			Message:   chat,
			StaffRole: user.StaffRole,
			Timestamp: time.Now(),
			GameID:    game.ID,
		})

		game.GamePublic.GeneralGameSettings.Mutex.Unlock()
		GameMapMutex.Lock()
		GameMap[game.GamePublic.ID] = *game
		GameMapMutex.Unlock()

		fmt.Println("Added new game chat")
		game.GamePublic.GeneralGameSettings.Mutex.RLock()
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.ID, "playerChatUpdate", game.GamePublic.Chats[len(game.GamePublic.Chats)-1])
		game.GamePublic.GeneralGameSettings.Mutex.RUnlock()
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
			UserID:    user.UserID,
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

	if game.GameCreatorBlacklist != nil {
		for _, id := range game.GameCreatorBlacklist {
			if id == user.UserID {
				socket.Emit("gameJoinStatusUpdate", bson.M{
					"status": "blacklisted",
				})

				return
			}
		}
	}

	game.GamePublic.GeneralGameSettings.Players = append(game.GamePublic.GeneralGameSettings.Players, *user)
	game.GamePublic.PublicPlayersState = append(game.GamePublic.PublicPlayersState, types.PlayerState{
		UserID:                   user.UserID,
		Connected:                true,
		LeftGame:                 false,
		CustomCardback:           "",
		CustomCardbackID:         "",
		IsLoader:                 false,
		IsRemakeVoting:           false,
		PingTime:                 0,
		Username:                 user.Username,
		PreviousGovernmentStatus: "",
		GovernmentStatus:         "",
	})

	game.GamePublic.GeneralGameSettings.Map[user.UserID] = game.GamePublic.PlayerCount
	game.GamePublic.PlayerCount = len(game.GamePublic.PublicPlayersState)
	game.GamePublic.GeneralGameSettings.Status = DisplayWaitingForPlayers(&game.GamePublic)

	GameMapMutex.Lock()
	GameMap[game.GamePublic.ID] = game
	GameMapMutex.Unlock()

	GameMapMutex.RLock()
	database.RedisDB.Set(ctx, "gamesMap", GameMap, 0)
	GameMapMutex.RUnlock()

	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	IO.JoinRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID, socket)
	socket.Emit("updateSeatForUser")
	IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID, "gameUpdate", game.GamePublic)
	// socket.Emit("gameUpdate", game.GamePublic)
}
