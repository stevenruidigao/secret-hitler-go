package socket

import (
	// "fmt"
	"secrethitler.io/types"
	"secrethitler.io/utils"
	// "strconv"
	// "time"
)

func SendInProgressGameUpdate(game *types.GamePrivate) {
	for i := range game.SeatedPlayers {
		gamePrivate := types.GamePrivate{
			GamePublic:   game.GamePublic,
			PlayerStates: game.SeatedPlayers[i].PlayerStates,
			// SeatedPlayers: make([]types.PlayerState, len(game.SeatedPlayers)),
		}

		// copy(gamePrivate.SeatedPlayers, game.GamePublic.PublicPlayerStates)
		// gamePrivate.SeatedPlayers[i] = game.SeatedPlayers[i]

		// if i == 0 {
		// a, _ := utils.MarshalJSON(gamePrivate)
		// fmt.Println("Seated Players (Game Update)", a, gamePrivate.SeatedPlayers[0])
		// }
		// game.SeatedPlayers[i].Socket.Emit("gameUpdate", gamePrivate)
		gamePrivate.GamePublic.Chats = append(game.SeatedPlayers[i].GameChats, game.GamePublic.Chats...)
		// fmt.Println("Merged Chats", gamePrivate.GamePublic.Chats)
		IO.BroadcastToRoom("/", "game-"+gamePrivate.GamePublic.ID+"-"+game.SeatedPlayers[i].ID, "gameUpdate", gamePrivate)
		// IO.BroadcastToRoom("/", "game-"+gamePrivate.GamePublic.ID+"-"+game.SeatedPlayers[i].ID, "gameUpdate2", gamePrivate)
	}

	gamePublic := types.GamePrivate{
		GamePublic: game.GamePublic,
	}

	gamePublic.GamePublic.Chats = append(game.UnseatedGameChats, game.GamePublic.Chats...)
	IO.BroadcastToRoom("/", "game-"+gamePublic.GamePublic.ID+"-observer", "gameUpdate", gamePublic)
}

func ShufflePolicies(game *types.GamePrivate) {
	game.Policies = make([]types.Policy, game.GamePublic.CustomGameSettings.DeckState.Fascist-game.GamePublic.TrackState.FascistPolicyCount+game.GamePublic.CustomGameSettings.DeckState.Liberal-game.GamePublic.TrackState.LiberalPolicyCount)

	for i := 0; i < game.GamePublic.CustomGameSettings.DeckState.Fascist-game.GamePublic.TrackState.FascistPolicyCount; i++ {
		game.Policies[i] = types.Policy{
			Cardback: "fascist",
		}
	}

	for i := game.GamePublic.CustomGameSettings.DeckState.Fascist - game.GamePublic.TrackState.FascistPolicyCount; i < game.GamePublic.CustomGameSettings.DeckState.Fascist-game.GamePublic.TrackState.FascistPolicyCount+game.GamePublic.CustomGameSettings.DeckState.Liberal-game.GamePublic.TrackState.LiberalPolicyCount; i++ {
		game.Policies[i] = types.Policy{
			Cardback: "liberal",
		}
	}

	for i := 0; i < len(game.Policies)-1; i++ {
		j := utils.RandInt(uint32(i), uint32(len(game.Policies)))
		policy := game.Policies[i]
		game.Policies[i] = game.Policies[j]
		game.Policies[j] = policy
	}
}
