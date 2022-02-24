package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	// "strconv"
	// "time"
)

func SendInProgressGameUpdate(game *types.GamePrivate) {
	for i := range game.SeatedPlayers {
		gamePrivate := types.GamePrivate{
			GamePublic:    game.GamePublic,
			SeatedPlayers: make([]types.PlayerState, len(game.SeatedPlayers)),
		}

		copy(gamePrivate.SeatedPlayers, game.GamePublic.PublicPlayerStates)
		gamePrivate.SeatedPlayers[i] = game.SeatedPlayers[i]

		if i == 0 {
			a, _ := utils.MarshalJSON(gamePrivate)
			fmt.Println("Seated Players (Game Update)", a, gamePrivate.SeatedPlayers[0])
		}
		// game.SeatedPlayers[i].Socket.Emit("gameUpdate", gamePrivate)
		IO.BroadcastToRoom("/", "game-"+gamePrivate.GamePublic.ID+"-"+game.SeatedPlayers[i].ID, "gameUpdate", gamePrivate)
		IO.BroadcastToRoom("/", "game-"+gamePrivate.GamePublic.ID+"-"+game.SeatedPlayers[i].ID, "gameUpdate2", gamePrivate)
	}
}
