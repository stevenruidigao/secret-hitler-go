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
		if currentPlayerCount == playerCount {
			break

		} else if playerCount > currentPlayerCount {
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

func Countdown(game *types.GamePrivate, timer int) {
	if timer == 0 {
		BeginGame(game)
		return
	}

	game.GamePublic.GeneralGameSettings.Status = "Game starts in " + strconv.Itoa(timer) + " second"

	if timer != 1 {
		game.GamePublic.GeneralGameSettings.Status += "s"
	}

	IO.BroadcastToRoom("/", "game-"+game.GeneralGameSettings.ID, "gameUpdate", game.GamePublic)

	time.AfterFunc(1*time.Second, func() {
		Countdown(game, timer-1)
	})
}
