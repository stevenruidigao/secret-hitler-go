package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	// "fmt"
	"strconv"
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
		gamePrivate.CardFlingerState = game.SeatedPlayers[i].CardFlingerState
		// fmt.Println("Merged Chats", gamePrivate.GamePublic.Chats)
		IO.BroadcastToRoom("/", "game-"+gamePrivate.GamePublic.GeneralGameSettings.ID+"-"+game.SeatedPlayers[i].ID, "gameUpdate", gamePrivate)
		// IO.BroadcastToRoom("/", "game-"+gamePrivate.GamePublic.GeneralGameSettings.ID+"-"+game.SeatedPlayers[i].ID, "gameUpdate2", gamePrivate)
	}

	gamePublic := types.GamePrivate{
		GamePublic: game.GamePublic,
	}

	gamePublic.GamePublic.Chats = append(game.UnseatedGameChats, game.GamePublic.Chats...)
	IO.BroadcastToRoom("/", "game-"+gamePublic.GamePublic.GeneralGameSettings.ID+"-observer", "gameUpdate", gamePublic)
}

func SendInProgressModChatUpdate(game *types.GamePrivate, chat types.PlayerChat) {

}

func ShufflePolicies(game *types.GamePrivate, start bool) {
	if game == nil {
		return
	}

	if start {
		game.GamePublic.TrackState.FascistPolicyCount = game.GamePublic.CustomGameSettings.TrackState.Fascist

		for i := 0; i < game.GamePublic.CustomGameSettings.TrackState.Fascist; i++ {
			game.GamePublic.TrackState.EnactedPolicies = append(game.GamePublic.TrackState.EnactedPolicies, types.Policy{
				CardBack: "fascist",
				Flipped:  true,
				Position: "fascist" + strconv.Itoa(i),
			})
		}

		game.GamePublic.TrackState.LiberalPolicyCount = game.GamePublic.CustomGameSettings.TrackState.Liberal

		for i := 0; i < game.GamePublic.CustomGameSettings.TrackState.Liberal; i++ {
			game.GamePublic.TrackState.EnactedPolicies = append(game.GamePublic.TrackState.EnactedPolicies, types.Policy{
				CardBack: "liberal",
				Flipped:  true,
				Position: "liberal" + strconv.Itoa(i),
			})
		}
	}

	game.Policies = make([]string, game.GamePublic.CustomGameSettings.DeckState.Fascist-game.GamePublic.TrackState.FascistPolicyCount+game.GamePublic.CustomGameSettings.DeckState.Liberal-game.GamePublic.TrackState.LiberalPolicyCount)

	for i := 0; i < game.GamePublic.CustomGameSettings.DeckState.Fascist-game.GamePublic.TrackState.FascistPolicyCount; i++ {
		game.Policies[i] = "fascist"
	}

	for i := game.GamePublic.CustomGameSettings.DeckState.Fascist - game.GamePublic.TrackState.FascistPolicyCount; i < game.GamePublic.CustomGameSettings.DeckState.Fascist-game.GamePublic.TrackState.FascistPolicyCount+game.GamePublic.CustomGameSettings.DeckState.Liberal-game.GamePublic.TrackState.LiberalPolicyCount; i++ {
		game.Policies[i] = "liberal"
	}

	for i := 0; i < len(game.Policies)-1; i++ {
		j := utils.RandInt(uint32(i), uint32(len(game.Policies)))
		policy := game.Policies[i]
		game.Policies[i] = game.Policies[j]
		game.Policies[j] = policy
	}
}
