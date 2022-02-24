package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
)

func BeginGame(game *types.GamePrivate) {
	fmt.Println("Game is starting...")

	game.GamePublic.GeneralGameSettings.TimeStarted = time.Now()
	customGameSettings := game.GamePublic.CustomGameSettings

	if !game.GamePublic.CustomGameSettings.Enabled {
		game.GamePublic.CustomGameSettings.HitlerZone = 3
		game.GamePublic.CustomGameSettings.VetoZone = 5
		game.GamePublic.CustomGameSettings.TrackState = types.CustomGameCounter{
			Liberal: 0,
			Fascist: 0,
		}
		game.GamePublic.CustomGameSettings.DeckState = types.CustomGameCounter{
			Liberal: 6,
			Fascist: 11,
		}

		if 5 <= game.GamePublic.PlayerCount && game.GamePublic.PlayerCount <= 6 {
			game.CustomGameSettings.FascistCount = 1
			game.CustomGameSettings.HitlerKnowsFascists = true
			game.CustomGameSettings.Powers = []string{"", "", "deckpeek", "bullet", "bullet"}

		} else if 7 <= game.GamePublic.PlayerCount && game.GamePublic.PlayerCount <= 8 {
			game.CustomGameSettings.FascistCount = 2
			game.CustomGameSettings.HitlerKnowsFascists = false
			game.CustomGameSettings.Powers = []string{"", "investigate", "election", "bullet", "bullet"}

		} else if 9 <= game.GamePublic.PlayerCount && game.GamePublic.PlayerCount <= 10 {
			game.CustomGameSettings.FascistCount = 3
			game.CustomGameSettings.HitlerKnowsFascists = false
			game.CustomGameSettings.Powers = []string{"investigate", "investigate", "election", "bullet", "bullet"}
		}

		game.GamePublic.CustomGameSettings.LiberalCount = game.GamePublic.PlayerCount - customGameSettings.FascistCount - 1
	}

	fmt.Println("Custom Game Settings", customGameSettings, game.GamePublic.CustomGameSettings)

	roles := make([]types.Role, game.GamePublic.PlayerCount)

	roles[0] = types.Role{
		CardName: "hitler",
		Icon:     1,
		Team:     "fascist",
	}

	for i := 0; i < game.GamePublic.CustomGameSettings.LiberalCount; i++ {
		roles[i+1] = types.Role{
			CardName: "liberal",
			Icon:     i % 6,
			Team:     "liberal",
		}
	}

	for i := 6; i < 6+game.GamePublic.CustomGameSettings.FascistCount; i++ {
		roles[i+game.GamePublic.CustomGameSettings.LiberalCount-6] = types.Role{
			CardName: "fascist",
			Icon:     i,
			Team:     "fascist",
		}
	}

	fmt.Println("Roles", roles)

	for i := 0; i < game.GamePublic.PlayerCount-1; i++ {
		j := utils.RandInt(uint32(i), uint32(game.GamePublic.PlayerCount))
		role := roles[i]
		roles[i] = roles[j]
		roles[j] = role
	}

	fmt.Println("Shuffled Roles", roles)

	game.GamePublic.GeneralGameSettings.Status = "Dealing roles..."

	for i := range game.GamePublic.PublicPlayerStates {
		game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = true
	}

	liberalPlayers := make([]types.PlayerState, game.GamePublic.CustomGameSettings.LiberalCount)
	liberalPlayerCount := 0
	fascistPlayers := make([]types.PlayerState, game.GamePublic.CustomGameSettings.FascistCount+1)
	fascistPlayerCount := 0

	for i := range game.SeatedPlayers {
		game.SeatedPlayers[i].Role = roles[i]
		game.SeatedPlayers[i].PlayerStates = make([]types.PlayerState, game.GamePublic.PlayerCount)

		for j := 0; j < game.GamePublic.PlayerCount; j++ {
			game.SeatedPlayers[i].PlayerStates[j] = types.PlayerState{
				NotificationStatus: "",
				NameStatus:         "",
			}
		}

		game.SeatedPlayers[i].PlayerStates[i].CardStatus.CardName = game.SeatedPlayers[i].Role.CardName

		if !game.GamePublic.GeneralGameSettings.DisableGamechat {
			game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.GameChats{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat: []types.GameChat{
					types.GameChat{
						Text: "The game begins and you receive the ",
					},
					types.GameChat{
						Text: game.SeatedPlayers[i].Role.CardName,
						Type: game.SeatedPlayers[i].Role.CardName,
					},
					types.GameChat{
						Text: " role and take seat #",
					},
					types.GameChat{
						Text: strconv.Itoa(i + 1),
						Type: "player",
					},
				},
			})

		} else {
			game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.GameChats{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat: []types.GameChat{
					types.GameChat{
						Text: "The game begins.",
					},
				},
			})
		}

		if game.SeatedPlayers[i].Role.CardName == "liberal" {
			liberalPlayers[liberalPlayerCount] = game.SeatedPlayers[i]
			liberalPlayerCount++

		} else {
			fascistPlayers[fascistPlayerCount] = game.SeatedPlayers[i]
			fascistPlayerCount++
		}
	}

	liberalElo := types.TeamElo{
		Overall:  1600,
		Seasonal: 1600,
	}

	fascistElo := types.TeamElo{
		Overall:  1600,
		Seasonal: 1600,
	}

	for i := range liberalPlayers {
		liberalElo.Overall += liberalPlayers[i].UserPublic.EloOverall
		liberalElo.Seasonal += liberalPlayers[i].UserPublic.EloSeason
	}

	for i := range fascistPlayers {
		fascistElo.Overall += fascistPlayers[i].UserPublic.EloOverall
		fascistElo.Seasonal += fascistPlayers[i].UserPublic.EloSeason
	}

	liberalElo.Overall /= float64(liberalPlayerCount)
	liberalElo.Seasonal /= float64(liberalPlayerCount)
	fascistElo.Overall /= float64(fascistPlayerCount)
	fascistElo.Seasonal /= float64(fascistPlayerCount)

	// fmt.Println("Seated Players (Start Game)", game.SeatedPlayers)
	SendInProgressGameUpdate(game)
}

func StartGame(game *types.GamePrivate) {
	game.GamePublic.GameState.TracksFlipped = true

	for i := 0; i < game.PlayerCount; i++ {
		j := utils.RandInt(uint32(i), uint32(game.PlayerCount))
		player := game.GamePublic.PublicPlayerStates[i]
		game.GamePublic.PublicPlayerStates[i] = game.GamePublic.PublicPlayerStates[j]
		game.GamePublic.PublicPlayerStates[j] = player
		game.GamePublic.PlayerMap[game.GamePublic.PublicPlayerStates[i].ID] = int(i)
		game.GamePublic.PlayerMap[game.GamePublic.PublicPlayerStates[j].ID] = int(j)
	}

	game.SeatedPlayers = make([]types.PlayerState, len(game.GamePublic.PublicPlayerStates))
	copy(game.SeatedPlayers, game.GamePublic.PublicPlayerStates)
	// fmt.Println("Players", game.SeatedPlayers, game.GamePublic.PublicPlayerStates)
	Countdown(game, 5)
}
