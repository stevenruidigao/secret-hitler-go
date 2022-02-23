package socket

import (
	"secrethitler.io/types"

	"strconv"
	"time"
)

func DisplayWaitingForPlayers(game *types.GamePrivate) string {
	status := ""
	currentPlayerCount := len(game.GamePublic.PublicPlayerStates)

	for _, playerCount := range game.GamePublic.PlayerCounts {
		if playerCount > currentPlayerCount {
			difference := playerCount - currentPlayerCount
			status = "Waiting for " + strconv.Itoa(difference) + " more player"

			if difference != 1 {
				status += "s"
			}

			status += "..."
			return status
		}
	}

	if game.GamePublic.GameState.Started {
		return status
	}

	game.GamePublic.GameState.Started = true
	StartGame(game)
	return status
}

func Countdown(game *types.GamePublic, timer int) {
	if timer == 0 {
		return
	}

	game.GeneralGameSettings.Status = "Game starts in " + strconv.Itoa(timer) + " seconds"
	IO.BroadcastToRoom("/", "game-"+game.ID, "gameUpdate", game)

	time.AfterFunc(1*time.Second, func() {
		Countdown(game, timer-1)
	})
}
