package socket

import (
	"secrethitler.io/types"

	"time"
)

func startGame(game types.GamePrivate) {
	game.GamePublic.GeneralGameSettings.TimeStarted = time.Now()
	customGameSettings := game.GamePublic.CustomGameSettings

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
}
