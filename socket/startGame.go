package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"time"
)

func BeginGame(game *types.GamePublic) {
	game.GeneralGameSettings.TimeStarted = time.Now()
	customGameSettings := game.CustomGameSettings

	if !customGameSettings.Enabled {
		customGameSettings.HitlerZone = 3
		customGameSettings.VetoZone = 5
		customGameSettings.TrackState = types.CustomTrackState{
			Liberal: 0,
			Fascist: 0,
		}
		customGameSettings.DeckState = types.CustomDeckState{
			Liberal: 6,
			Fascist: 11,
		}
	}

	roles := make([]types.Role, game.PlayerCount)

	roles[0] = types.Role{
		CardName: "hitler",
		Icon:     1,
		Team:     "fascist",
	}

	for i := 0; i < game.CustomGameSettings.LiberalCount; i++ {
		roles[i+1] = types.Role{
			CardName: "liberal",
			Icon:     i % 6,
			Team:     "liberal",
		}
	}

	for i := 6; i < 6+game.CustomGameSettings.FascistCount; i++ {
		roles[i+game.CustomGameSettings.LiberalCount-5] = types.Role{
			CardName: "fascist",
			Icon:     i,
			Team:     "fascist",
		}
	}

	for i := 0; i < game.PlayerCount-1; i++ {
		j := utils.RandInt(uint32(i), uint32(game.PlayerCount))
		role := roles[i]
		roles[i] = roles[j]
		roles[j] = role
	}

}

func StartGame(game *types.GamePrivate) {
	game.GamePublic.GameState.TracksFlipped = true

	for i := 0; i < game.PlayerCount; i++ {
		j := utils.RandInt(uint32(i), uint32(game.PlayerCount))
		player := game.GamePublic.PublicPlayerStates[i]
		game.GamePublic.PublicPlayerStates[i] = game.PublicPlayerStates[j]
		game.GamePublic.PublicPlayerStates[j] = player
	}

	copy(game.GamePublic.PublicPlayerStates, game.SeatedPlayers)
	Countdown(&game.GamePublic, 20)
}
