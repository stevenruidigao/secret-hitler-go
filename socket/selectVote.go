package socket

import (
	"secrethitler.io/types"
	// "secrethitler.io/utils"

	"fmt"
	// "strconv"
	"time"
	// "github.com/googollee/go-socket.io"
)

func SelectVote(playerState *types.PlayerState, game *types.GamePrivate, vote bool, force bool) {
	if playerState == nil || game == nil {
		return
	}

	if game.GamePublic.GameState.Frozen && !force {
		fmt.Println("Frozen")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+playerState.UserPublic.ID, "sendAlert", "An AEM member has prevented this game from proceeding. Please wait.")
		return
	}

	if game.GamePublic.GeneralGameSettings.Remade && !force {
		fmt.Println("Remade")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+playerState.UserPublic.ID, "sendAlert", "This game has been remade and is now no longer playable.")
		return
	}

	if game.GeneralGameSettings.Tourny && game.GeneralGameSettings.TournyInfo.Cancelled {
		return
	}

	votedPlayerCount := 0

	for i := range game.SeatedPlayers {
		if game.SeatedPlayers[i].VoteStatus.Voted {
			votedPlayerCount++
		}
	}

	fmt.Println("Condition", len(game.SeatedPlayers), votedPlayerCount)

	if len(game.SeatedPlayers) != votedPlayerCount {
		if !playerState.VoteStatus.Voted || force {
			playerState.VoteStatus.Voted = true

		} else {
			if playerState.VoteStatus.VotedYes {
				playerState.VoteStatus.Voted = !vote

			} else {
				playerState.VoteStatus.Voted = vote
			}
		}

		playerState.VoteStatus.VotedYes = playerState.VoteStatus.Voted && vote
		game.GamePublic.PublicPlayerStates[playerState.Index].Loader = !playerState.VoteStatus.Voted
		notificationStatus := "notification"

		if playerState.VoteStatus.Voted {
			notificationStatus = "selected"
		}

		if vote {
			playerState.CardFlingerState = []types.CardFlinger{
				types.CardFlinger{
					Position:           "middle-left",
					NotificationStatus: notificationStatus,
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   true,
						CardFront: "ballot",
						CardBack:  "ja",
					},
				},
				types.CardFlinger{
					Position:           "middle-right",
					NotificationStatus: "notification",
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   true,
						CardFront: "ballot",
						CardBack:  "nein",
					},
				},
			}

		} else {
			playerState.CardFlingerState = []types.CardFlinger{
				types.CardFlinger{
					Position:           "middle-left",
					NotificationStatus: "notification",
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   true,
						CardFront: "ballot",
						CardBack:  "ja",
					},
				},
				types.CardFlinger{
					Position:           "middle-right",
					NotificationStatus: notificationStatus,
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   true,
						CardFront: "ballot",
						CardBack:  "nein",
					},
				},
			}
		}

		SendInProgressGameUpdate(game)
		votedPlayerCount = 0

		for i := range game.SeatedPlayers {
			if game.SeatedPlayers[i].VoteStatus.Voted && !game.SeatedPlayers[i].Dead {
				votedPlayerCount++
			}
		}

		fmt.Println("Voted:", votedPlayerCount)
		fmt.Println("Living", game.GamePublic.GeneralGameSettings.LivingPlayerCount)

		if votedPlayerCount == game.GamePublic.GeneralGameSettings.LivingPlayerCount {
			game.GamePublic.GeneralGameSettings.Status = "Tallying results of ballots..."

			for i := range game.SeatedPlayers {
				if len(game.SeatedPlayers[i].CardFlingerState) > 0 {
					game.SeatedPlayers[i].CardFlingerState[0].Action = ""
					game.SeatedPlayers[i].CardFlingerState[1].Action = ""
					game.SeatedPlayers[i].CardFlingerState[0].CardStatus.Flipped = false
					game.SeatedPlayers[i].CardFlingerState[1].CardStatus.Flipped = false
				}
			}

			SendInProgressGameUpdate(game)

			time.AfterFunc(200*time.Millisecond, func() {
				for i := range game.SeatedPlayers {
					game.SeatedPlayers[i].CardFlingerState = []types.CardFlinger{}
				}

				SendInProgressGameUpdate(game)
			})

			time.AfterFunc(2500*time.Millisecond, func() {
				if game.GamePublic.GeneralGameSettings.Timer > 0 && game.Timer != nil {
					game.Timer.Stop()
					game.GamePublic.GameState.TimedMode = false
				}

				FlipBallotCards(game)
			})
		}
	}
}
