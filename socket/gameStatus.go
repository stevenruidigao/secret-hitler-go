package socket

import (
	"secrethitler.io/types"

	"strconv"
	"time"
)

func DisplayWaitingForPlayers(game *types.GamePublic) string {
	status := ""
	currentPlayerCount := len(game.PublicPlayersState)

	for _, playerCount := range game.PlayerCounts {
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

	Countdown(game, 20)
	return status
}

func Countdown(game *types.GamePublic, timer int) {
	game.GeneralGameSettings.Status = "Game starts in " + strconv.Itoa(timer) + " seconds"
	IO.BroadcastToRoom("/", "game-"+game.ID, "gameUpdate", game)

	time.AfterFunc(1*time.Second, func() {
		Countdown(game, timer-1)
	})
}
