package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
)

func FailedElection(game *types.GamePrivate) {
	game.GamePublic.TrackState.ElectionTrackerCount++

	if game.GamePublic.TrackState.ElectionTrackerCount >= 3 {
		game.GameState.PreviousElectedGovernment = []int{}

		if !game.GamePublic.GeneralGameSettings.DisableGamechat {
			chat := types.PlayerChat{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat: []types.GameChat{
					types.GameChat{
						Text: "The third consecutive election has failed and the top policy is enacted.",
					},
				},
			}

			for i := range game.SeatedPlayers {
				game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, chat)
			}

			game.UnseatedGameChats = append(game.UnseatedGameChats, chat)
		}

		if game.GamePublic.GameState.UndrawnPolicyCount == 0 {
			ShufflePolicies(game, false)
		}

		game.GamePublic.GameState.UndrawnPolicyCount--

		time.AfterFunc(500*time.Millisecond, func() {
			policy := game.Policies[0]
			game.Policies = game.Policies[1:]
			EnactPolicy(game, policy)
		})

	} else {
		if game.GamePublic.GeneralGameSettings.Timer > 0 {
			if game.Timer != nil {
				game.Timer.Stop()
				game.Timer = nil
			}

			game.GamePublic.GameState.TimedMode = true

			time.AfterFunc(500*time.Millisecond, func() {
				StartElection(game, -1)
			})

			game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
				if game.GamePublic.GameState.TimedMode && game.GamePublic.GameState.Phase == "selectingChancellor" {
					clickActionInfo, ok := game.GamePublic.GameState.ClickActionInfo[1].([]int)

					if !ok {
						return
					}

					chancellorIndex := clickActionInfo[utils.RandInt(0, uint32(len(clickActionInfo)))]
					game.GamePublic.GameState.PendingChancellorIndex = -1
					game.GamePublic.GameState.TimedMode = false

					SelectChancellor(&game.SeatedPlayers[game.GameState.PresidentIndex].UserPublic, game, chancellorIndex, false)
				}
			})
		}
	}
}

func PassedElection(game *types.GamePrivate) {
	game.GamePublic.GameState.ChancellorIndex = game.GameState.PendingChancellorIndex
	fmt.Println("Passed election")

	if game.GamePublic.GameState.PreviousElectedGovernment[0] != -1 {
		// game.SeatedPlayers[game.GameState.PreviousElectedGovernment[0]].PlayerStates[game.GameState.PreviousElectedGovernment[0]].Claim = ""
		// game.SeatedPlayers[game.GameState.PreviousElectedGovernment[1]].PlayerStates[game.GameState.PreviousElectedGovernment[1]].Claim = ""

		for i := range game.GamePublic.GameState.PreviousElectedGovernment {
			game.SeatedPlayers[game.GameState.PreviousElectedGovernment[i]].PlayerStates[game.GameState.PreviousElectedGovernment[i]].Claim = ""
			IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+game.SeatedPlayers[game.GameState.PreviousElectedGovernment[i]].UserPublic.ID, "removeClaim")
		}
	}

	game.GamePublic.GeneralGameSettings.Status = "Waiting on presidential discard."
	game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].Loader = true

	if !game.GamePublic.GeneralGameSettings.Experienced && !game.GamePublic.GeneralGameSettings.DisableGamechat {
		game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].GameChats = append(game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].GameChats, types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: "As president, you must select one policy to discard.",
				},
			},
		})
	}

	if game.GamePublic.GameState.UndrawnPolicyCount < 3 {
		ShufflePolicies(game, false)
	}

	// policies := game.Policies[:3]
	game.CurrentElectionPolicies = make([]string, 3)
	copy(game.CurrentElectionPolicies, game.Policies[:3])
	game.Policies = game.Policies[3:]
	validHand := true
	comment := "has just received an invalid hand!\n"

	// for i := range policies {
	// 	game.CurrentElectionPolicies[i] = policies[i]
	// }

	for i := range game.CurrentElectionPolicies {
		comment += game.CurrentElectionPolicies[i]

		if i != len(game.CurrentElectionPolicies)-1 {
			comment += ", "
		}

		if game.CurrentElectionPolicies[i] != "fascist" && game.CurrentElectionPolicies[i] != "liberal" {
			validHand = false
		}
	}

	if !validHand {
		gameType := "Ranked"

		if game.GamePublic.GeneralGameSettings.Casual {
			gameType = "Casual"

		} else if game.GamePublic.GeneralGameSettings.Practice {
			gameType = "Practice"
		}

		MakeReport(game, types.Report{
			ReportedPlayerID:       game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.ID,
			ReportedPlayerSeat:     game.GamePublic.GameState.PresidentIndex + 1,
			ReportedPlayerRole:     game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].Role.CardName,
			Comment:                comment,
			GameElectionCount:      game.GamePublic.GeneralGameSettings.ElectionCount,
			GameName:               game.GamePublic.GeneralGameSettings.Name,
			GameID:                 game.GamePublic.GeneralGameSettings.ID,
			GameType:               gameType,
			ReportedPlayerUsername: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.ID,
			ReportType:             "report",
		})
	}

	hiddenInfoChat := []types.GameChat{
		types.GameChat{
			Text: "President",
		},
		types.GameChat{
			Text: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PresidentIndex+1) + "}",
			Type: "player",
		},
		types.GameChat{
			Text: " received ",
		},
	}

	for i := range game.CurrentElectionPolicies {
		if game.CurrentElectionPolicies[i] == "fascist" {
			hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
				Text: "R",
				Type: "fascist",
			})

		} else {
			hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
				Text: "B",
				Type: "liberal",
			})
		}
	}

	hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
		Text: ".",
	})

	modChat := types.PlayerChat{
		Timestamp: time.Now(),
		GameChat:  true,
		Chat:      hiddenInfoChat,
	}

	game.HiddenInfoChat = append(game.HiddenInfoChat, modChat)
	SendInProgressModChatUpdate(game, modChat)

	game.Summary.Logs = append(game.Summary.Logs, struct {
		PresidentHand []string `bson:"presidentHand" json:"presidentHand"`
	}{
		PresidentHand: game.CurrentElectionPolicies,
	})

	game.CardFlingerState = []types.CardFlinger{
		types.CardFlinger{
			Position: "middle-far-left",
			Action:   "active",
			CardStatus: types.CardStatus{
				Flipped:   true,
				CardFront: "policy",
				CardBack:  game.CurrentElectionPolicies[0] + "p",
			},
			Discard: true,
		},
		types.CardFlinger{
			Position: "middle-center",
			Action:   "active",
			CardStatus: types.CardStatus{
				Flipped:   true,
				CardFront: "policy",
				CardBack:  game.CurrentElectionPolicies[1] + "p",
			},
			Discard: true,
		},
		types.CardFlinger{
			Position: "middle-far-right",
			Action:   "active",
			CardStatus: types.CardStatus{
				Flipped:   true,
				CardFront: "policy",
				CardBack:  game.CurrentElectionPolicies[2] + "p",
			},
			Discard: true,
		},
	}

	game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState = game.CardFlingerState
	game.GamePublic.GameState.UndrawnPolicyCount--
	SendInProgressGameUpdate(game)

	time.AfterFunc(200*time.Millisecond, func() {
		game.GamePublic.GameState.UndrawnPolicyCount--
		SendInProgressGameUpdate(game)
	})

	time.AfterFunc(400*time.Millisecond, func() {
		game.GamePublic.GameState.UndrawnPolicyCount--
		SendInProgressGameUpdate(game)
	})

	time.AfterFunc(200*time.Millisecond, func() {
		for i := range game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState {
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState[i].CardStatus.Flipped = true
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState[i].NotificationStatus = "notification"
		}

		game.GamePublic.GameState.Phase = "presidentSelectingPolicy"
		game.GamePublic.GameState.PreviousElectedGovernment = []int{game.GamePublic.GameState.PresidentIndex, game.GamePublic.GameState.ChancellorIndex}

		if game.GamePublic.GeneralGameSettings.Timer > 0 {
			if game.Timer != nil {
				game.Timer.Stop()
				game.Timer = nil
			}

			game.GamePublic.GameState.TimedMode = true

			game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
				if game.GameState.TimedMode {
					game.GameState.TimedMode = false
					SelectPresidentPolicy(&game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic, game, int(utils.RandInt(0, uint32(len(game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState)))), true)

					game.ReplayGameChats = append(game.ReplayGameChats, types.PlayerChat{
						Timestamp: time.Now(),
						GameChat:  true,
						Chat: []types.GameChat{
							types.GameChat{
								Text: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PresidentIndex+1) + "}",
								Type: "player",
							},
							types.GameChat{
								Text: " was forced by the timer to select a random policy to discard.",
							},
						},
					})
				}
			})
		}
	})
}
