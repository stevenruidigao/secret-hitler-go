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

		if !user.GameSettings.Incognito || isAEM {
			viewableList = append(viewableList, types.UserPublic{
				UserID:        user.UserID,
				Username:      user.Username,
				Created:       user.Created,
				Wins:          user.Wins,
				Losses:        user.Losses,
				RainbowWins:   user.RainbowWins,
				RainbowLosses: user.RainbowLosses,
				GameSettings: types.GameSettings{
					Private:           user.GameSettings.Private,
					DisableVisibleElo: user.GameSettings.DisableVisibleElo,
					DisableStaffColor: user.GameSettings.DisableStaffColor,
					// Blacklists are sent in the sendUserGameSettings event.
					// Blacklist: user.GameSettings.Blacklist,
					Cardback:                user.GameSettings.Cardback,
					CardbackID:              user.GameSettings.CardbackID,
					TournyWins:              user.GameSettings.TournyWins,
					SeasonAwards:            user.GameSettings.SeasonAwards,
					SpecialTournamentStatus: user.GameSettings.SpecialTournamentStatus,
					Incognito:               user.GameSettings.Incognito,
				},
				EloOverall:          math.Floor(user.EloOverall),
				EloSeason:           math.Floor(user.EloSeason),
				Status:              user.Status,
				Seasons:             user.Seasons,
				TimeLastGameCreated: user.TimeLastGameCreated,
				StaffRole:           user.StaffRole,
				Contributor:         user.Contributor,
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

func UpdateUserStatus(user types.UserPublic, game *types.GamePublic, override string) {
	UserMapMutex.RLock()

	for key, _ := range UserMap {
		if UserMap[key].UserID == user.UserID {
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
		playerNumber, ok := game.GamePublic.GeneralGameSettings.Map[user.UserID].(int)

		fmt.Println("Player number:", playerNumber)

		if ok {
			game.GamePublic.GeneralGameSettings.Players[playerNumber].LeftGame = false
			game.GamePublic.GeneralGameSettings.Players[playerNumber].Connected = true

			if game.GamePublic.GeneralGameSettings.TimeAbandoned != nil {
				game.GamePublic.GeneralGameSettings.TimeAbandoned = nil
			}

			socket.Emit("updateSeatForUser", true)
			UpdateUserStatus(*user, &game.GamePublic, "playing")

		} else {
			UpdateUserStatus(*user, &game.GamePublic, "observing")
		}
	}

	fmt.Println("Updated user status")

	IO.JoinRoom("/", "game-"+id, socket)
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
