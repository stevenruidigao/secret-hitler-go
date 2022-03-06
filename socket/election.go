package socket

import (
	"secrethitler.io/types"
	// "secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
	// "github.com/googollee/go-socket.io"
)

func FlipBallotCards(game *types.GamePrivate) {
	consensus := true
	vote := false

	for i := range game.GamePublic.PublicPlayerStates {
		if !game.GamePublic.PublicPlayerStates[i].Dead {
			vote = game.GamePublic.PublicPlayerStates[i].VoteStatus.VotedYes
		}
	}

	for i := range game.GamePublic.PublicPlayerStates {
		if game.GamePublic.PublicPlayerStates[i].VoteStatus.Voted != vote && !game.GamePublic.PublicPlayerStates[i].Dead {
			consensus = false
		}
	}

	votes := make([]bool, len(game.GamePublic.PublicPlayerStates))

	for i := range game.GamePublic.PublicPlayerStates {
		if !game.GamePublic.PublicPlayerStates[i].Dead {
			cardBack := types.CardBack{}

			if game.SeatedPlayers[i].VoteStatus.VotedYes {
				cardBack.CardName = "ja"

			} else {
				cardBack.CardName = "nein"
			}

			game.GamePublic.PublicPlayerStates[i].CardStatus.CardBack = cardBack
			game.GamePublic.PublicPlayerStates[i].CardStatus.Flipped = true
			votes[i] = game.SeatedPlayers[i].VoteStatus.VotedYes
			game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = true
		}
	}

	game.Summary.Logs = append(game.Summary.Logs, struct {
		Votes []bool `bson:"votes" json:"votes"`
	}{
		Votes: votes,
	})

	SendInProgressGameUpdate(game)

	timeout := 6000 * time.Millisecond

	if consensus {
		timeout = 1500 * time.Millisecond
	}

	time.AfterFunc(timeout, func() {
		chat := types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
		}

		for i := range game.GamePublic.PublicPlayerStates {
			game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = false
		}

		time.AfterFunc(500*time.Millisecond, func() {
			for i := range game.GamePublic.PublicPlayerStates {
				game.GamePublic.PublicPlayerStates[i].CardStatus.Flipped = false
			}

			SendInProgressGameUpdate(game)
		})

		yesVotes := 0

		for i := range game.SeatedPlayers {
			if game.SeatedPlayers[i].VoteStatus.VotedYes {
				yesVotes++
			}
		}

		fmt.Println("Yes votes:", yesVotes)

		if yesVotes > game.GamePublic.GeneralGameSettings.LivingPlayerCount/2 {
			game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].GovernmentStatus = "isPresident"
			game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PendingChancellorIndex].GovernmentStatus = "isChancellor"

			chat.Chat = []types.GameChat{
				types.GameChat{
					Text: "The election passes.",
				},
			}

			if !game.GamePublic.GeneralGameSettings.Experienced && !game.GamePublic.GeneralGameSettings.DisableGamechat {
				for i := range game.SeatedPlayers {
					game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, chat)
				}

				game.UnseatedGameChats = append(game.UnseatedGameChats, chat)
			}

			if game.GamePublic.TrackState.FascistPolicyCount >= game.GamePublic.CustomGameSettings.HitlerZone && game.SeatedPlayers[game.GamePublic.GameState.PendingChancellorIndex].Role.CardName == "hitler" {
				numberText := ""

				if game.GamePublic.CustomGameSettings.HitlerZone%10 == 1 && game.GamePublic.CustomGameSettings.HitlerZone != 11 {
					numberText = strconv.Itoa(game.GamePublic.CustomGameSettings.HitlerZone) + "st"

				} else if game.GamePublic.CustomGameSettings.HitlerZone%10 == 2 && game.GamePublic.CustomGameSettings.HitlerZone != 12 {
					numberText = strconv.Itoa(game.GamePublic.CustomGameSettings.HitlerZone) + "nd"

				} else if game.GamePublic.CustomGameSettings.HitlerZone%10 == 3 && game.GamePublic.CustomGameSettings.HitlerZone != 13 {
					numberText = strconv.Itoa(game.GamePublic.CustomGameSettings.HitlerZone) + "rd"

				} else {
					numberText = strconv.Itoa(game.GamePublic.CustomGameSettings.HitlerZone) + "th"
				}

				time.AfterFunc(1000*time.Millisecond, func() {
					for i := range game.GamePublic.PublicPlayerStates {
						game.GamePublic.PublicPlayerStates[i].CardStatus.CardFront = "secretrole"
						game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = true
						game.GamePublic.PublicPlayerStates[i].CardStatus.CardBack = game.SeatedPlayers[i].Role
					}

					if !game.GamePublic.GeneralGameSettings.DisableGamechat {
						chat := types.PlayerChat{
							Timestamp: time.Now(),
							GameChat:  true,
							Chat: []types.GameChat{
								types.GameChat{
									Text: "Hitler",
									Type: "hitler",
								},
								types.GameChat{
									Text: " has been elected chancellor after the " + numberText + " fascist policy has been enacted.",
								},
							},
						}

						for i := range game.SeatedPlayers {
							game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, chat)
						}

						game.UnseatedGameChats = append(game.UnseatedGameChats, chat)
					}

					game.GamePublic.GameState.AudioCue = "fascistsWinHitlerElected"
					SendInProgressGameUpdate(game)
				})

				time.AfterFunc(2*time.Second, func() {
					game.GamePublic.GameState.AudioCue = ""

					for i := range game.GamePublic.PublicPlayerStates {
						game.GamePublic.PublicPlayerStates[i].CardStatus.Flipped = true
					}

					fmt.Println("Game Complete")
					CompleteGame(game, "fascist")
				})

			} else {
				fmt.Println("Election Passed")
				PassedElection(game)
			}

		} else {
			if !game.GamePublic.GeneralGameSettings.DisableGamechat {
				chat.Chat = []types.GameChat{
					types.GameChat{
						Text: "The election fails and the election tracker moves forward. (" + strconv.Itoa(game.GamePublic.TrackState.ElectionTrackerCount+1) + ")",
					},
				}

				for i := range game.SeatedPlayers {
					game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, chat)
				}

				game.UnseatedGameChats = append(game.UnseatedGameChats, chat)
				game.GamePublic.GameState.PendingChancellorIndex = -1
			}

			fmt.Println("Election Failed")
			FailedElection(game)
		}

		SendInProgressGameUpdate(game)
	})
}
