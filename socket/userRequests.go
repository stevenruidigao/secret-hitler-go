package socket

import (
	// "secrethitler.io/constants"
	"secrethitler.io/types"
	// "secrethitler.io/utils"

	"fmt"
	"math"

	"github.com/googollee/go-socket.io"
)

func GetGameList(isAEM bool) []types.GeneralGameSettings {
	viewableList := []types.GeneralGameSettings{}
	GameMapMutex.RLock()

	for _, game := range GameMap {
		if !game.GamePublic.GeneralGameSettings.Private || isAEM {
			viewableList = append(viewableList, game.GamePublic.GeneralGameSettings)

			if !isAEM {
				viewableList[len(viewableList)-1].GameCreatorName = ""
			}
		}
	}

	GameMapMutex.RUnlock()
	return viewableList
}

func GetUserList(isAEM bool) interface{} {
	viewableList := []types.UserPublic{}
	UserMapMutex.RLock()

	for _, user := range UserMap {
		// fmt.Println(user)

		if user == nil {
			continue
		}

		if !user.GameSettings.Incognito || isAEM {
			viewableList = append(viewableList, types.UserPublic{
				ID:            user.UserPublic.ID,
				Username:      user.UserPublic.Username,
				Created:       user.UserPublic.Created,
				Wins:          user.UserPublic.Wins,
				Losses:        user.UserPublic.Losses,
				RainbowWins:   user.UserPublic.RainbowWins,
				RainbowLosses: user.UserPublic.RainbowLosses,
				// GameSettings: types.GameSettings{
				// 	Private:           user.GameSettings.Private,
				// 	DisableVisibleElo: user.GameSettings.DisableVisibleElo,
				// 	DisableStaffColor: user.GameSettings.DisableStaffColor,
				// 	// Blacklists are sent in the sendUserGameSettings event.
				// 	// Blacklist: user.GameSettings.Blacklist,
				// 	Cardback:                user.GameSettings.Cardback,
				// 	CardbackID:              user.GameSettings.CardbackID,
				// 	TournyWins:              user.GameSettings.TournyWins,
				// 	SeasonAwards:            user.GameSettings.SeasonAwards,
				// 	SpecialTournamentStatus: user.GameSettings.SpecialTournamentStatus,
				// 	Incognito:               user.GameSettings.Incognito,
				// },
				EloOverall:          math.Floor(user.UserPublic.EloOverall),
				EloSeason:           math.Floor(user.UserPublic.EloSeason),
				Status:              user.UserPublic.Status,
				Seasons:             user.UserPublic.Seasons,
				TimeLastGameCreated: user.UserPublic.TimeLastGameCreated,
				StaffRole:           user.UserPublic.StaffRole,
				Contributor:         user.UserPublic.Contributor,
			})
		}
	}

	UserMapMutex.RUnlock()

	userList := struct {
		List []types.UserPublic `bson:"list" json:"list"`
	}{List: viewableList}

	// fmt.Println(userList)

	return userList
}

func SendGameList(socket socketio.Conn) {
	user := GetUser(socket)
	socket.Emit("gameList", GetGameList(user != nil && user.StaffRole != "" && user.StaffRole != "altmod" && user.StaffRole != "veteran"))
}

func SendUserList(socket socketio.Conn) {
	user := GetUser(socket)
	socket.Emit("userList", GetUserList(user != nil && user.StaffRole != "" && user.StaffRole != "altmod" && user.StaffRole != "veteran"))
}

func UpdateUserStatus(user *types.UserPublic, game *types.GamePublic, override string) {
	UserMapMutex.RLock()

	for key := range UserMap {
		if UserMap[key].ID == user.ID {
			statusType := override

			if override == "" {
				statusType = "none"
			}

			if game != nil && !game.Unlisted {
				statusType = "playing"

				if game.GeneralGameSettings.Private {
					statusType = "private"

				} else if game.GeneralGameSettings.Rainbow {
					statusType = "rainbow"
				}
			}

			gameID := ""

			if game != nil && !game.Unlisted {
				gameID = game.GeneralGameSettings.ID
			}

			if user.Status == nil {
				user.Status = &types.UserStatus{
					Type:   statusType,
					GameID: gameID,
				}

				break
			}

			user.Status.Type = statusType
			user.Status.GameID = gameID
			break
		}
	}

	UserMapMutex.RUnlock()
}

func SendGameInfo(socket socketio.Conn, user *types.UserPublic, id string) {
	GameMapMutex.RLock()
	game := GameMap[id]
	GameMapMutex.RUnlock()

	if user != nil {
		playerNumber := game.GamePublic.PlayerMap[user.ID]

		fmt.Println("Player number:", playerNumber)

		if playerNumber > 0 {
			game.GamePublic.PublicPlayerStates[playerNumber-1].LeftGame = false
			game.GamePublic.PublicPlayerStates[playerNumber-1].Connected = true

			if game.GamePublic.GeneralGameSettings.TimeAbandoned != nil {
				game.GamePublic.GeneralGameSettings.TimeAbandoned = nil
			}

			socket.Emit("updateSeatForUser", true)
			UpdateUserStatus(user, &game.GamePublic, "playing")
			IO.JoinRoom("/", "game-"+id, socket)
			IO.JoinRoom("/", "game-"+game.GamePublic.ID+"-"+user.ID, socket)
			SendInProgressGameUpdate(game)
			socket.Emit("joinGameRedirect", id)
			return

		} else {
			UpdateUserStatus(user, &game.GamePublic, "observing")
			IO.JoinRoom("/", "game-"+game.GamePublic.ID+"-observer", socket)
		}
	}

	fmt.Println("Updated user status")

	IO.JoinRoom("/", "game-"+id, socket)
	// IO.JoinRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-observer", socket)
	socket.Emit("gameUpdate", game.GamePublic)
	socket.Emit("joinGameRedirect", id)
}

/*sendGameInfo = (socket, uid) => {
	const game = games[uid];
	const { passport } = socket.handshake.session;

	if (game) {
		if (passport && Object.keys(passport).length) {
			const player = game.publicPlayersState.find(player => player.userName === passport.user);

			if (player) {
				player.leftGame = false;
				player.connected = true;
				if (game.general) game.general.timeAbandoned = null;
				socket.emit('updateSeatForUser', true);
				updateUserStatus(passport, game);
			} else {
				updateUserStatus(passport, game, 'observing');
			}
		}

		socket.join(uid);
		sendInProgressGameUpdate(game);
		socket.emit('joinGameRedirect', game.general.uid);
	} else {
		Game.findOne({ uid }).then((game, err) => {
			if (err) {
				console.log(err, 'game err retrieving for replay');
			}

			socket.emit('manualReplayRequest', game ? game.uid : '');
		});
	}
};*/
