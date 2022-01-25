package socket

import (
//	"secrethitler.io/constants"
	"secrethitler.io/types"

	"math"

	"github.com/googollee/go-socket.io"
)

func SendUserList(socket socketio.Conn) {
	currentUser := GetUser(socket)

	userList := struct{
		List []types.User `bson:"list" json:"list"`
	}{}

	if currentUser == nil {
		socket.Emit("userList", userList)

		return
	}

	viewableList := []types.User{}

	if currentUser.StaffRole != "" && currentUser.StaffRole != "altmod" && currentUser.StaffRole != "veteran" {
		for _, user := range UserList {
			viewableList = append(viewableList, types.User {
				Username: user.Username,
				Wins: user.Wins,
				Losses: user.Losses,
				RainbowWins: user.RainbowWins,
				RainbowLosses: user.RainbowLosses,
				GameSettings: types.GameSettings {
					Private: user.GameSettings.Private,
					DisableVisibleElo: user.GameSettings.DisableVisibleElo,
					DisableStaffColor: user.GameSettings.DisableStaffColor,
					// Blacklists are sent in the sendUserGameSettings event.
					// Blacklist: user.GameSettings.Blacklist,
					Cardback: user.GameSettings.Cardback,
					CardbackID: user.GameSettings.CardbackID,
					TournyWins: user.GameSettings.TournyWins,
					SeasonAwards: user.GameSettings.SeasonAwards,
					SpecialTournamentStatus: user.GameSettings.SpecialTournamentStatus,
					Incognito: user.GameSettings.Incognito,
				},

				EloOverall: math.Floor(user.EloOverall),
				EloSeason: math.Floor(user.EloSeason),
				Status: user.Status,
				Seasons: user.Seasons,
				TimeLastGameCreated: user.TimeLastGameCreated,
				StaffRole: user.StaffRole,
				Contributor: user.Contributor,
			})
		}
	}
}
