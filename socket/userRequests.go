package socket

import (
	// "secrethitler.io/constants"
	"secrethitler.io/types"

	// "fmt"
	"math"

	"github.com/googollee/go-socket.io"
)

func GetGameList(isAEM bool) []types.GeneralGameSettings {
	viewableList := []types.GeneralGameSettings{}
	GameListMutex.RLock()

	for _, game := range GameList {
		if !game.Game.GeneralGameSettings.Private || isAEM {
			viewableList = append(viewableList, game.Game.GeneralGameSettings)
		}
	}

	GameListMutex.RUnlock()
	return viewableList
}

func GetUserList(isAEM bool) interface{} {
	viewableList := []types.User{}
	UserListMutex.RLock()

	for _, user := range UserList {
		if !user.GameSettings.Incognito || isAEM {
			viewableList = append(viewableList, types.User{
				Username:      user.Username,
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

	UserListMutex.RUnlock()

	userList := struct {
		List []types.User `bson:"list" json:"list"`
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

func UpdateUserStatus(user types.User, game types.Game) {
	UserListMutex.RLock()

	for i, _ := range UserList {
		if UserList[i].UserID == user.UserID {
			statusType := "none"

			if !game.Unlisted {
				statusType = "playing"

				if game.GeneralGameSettings.Private {
					statusType = "private"

				} else if game.GeneralGameSettings.Rainbow {
					statusType = "rainbow"
				}
			}

			gameID := ""

			if !game.Unlisted {
				gameID = game.GeneralGameSettings.ID
			}

			user.Status.Type = statusType
			user.Status.GameID = gameID
			break
		}
	}

	UserListMutex.RUnlock()
}
