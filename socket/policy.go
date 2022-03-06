package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
)

func EnactPolicy(game *types.GamePrivate, policy string) {
	if game.Lock.SelectChancellor {
		game.Lock.SelectChancellor = false
	}

	if game.Lock.SelectChancellorVoteOnVeto {
		game.Lock.SelectChancellorVoteOnVeto = false
	}

	if game.Lock.SelectChancellorPolicy {
		game.Lock.SelectChancellorPolicy = false
	}

	if game.Lock.PolicyPeek {
		game.Lock.PolicyPeek = false
	}

	if game.Lock.PolicyPeekAndDrop {
		game.Lock.PolicyPeekAndDrop = false
	}

	if game.Lock.SelectPlayerToExecute {
		game.Lock.SelectPlayerToExecute = false
	}

	if game.Lock.ExecutePlayer {
		game.Lock.ExecutePlayer = false
	}

	if game.Lock.SelectSpecialElection {
		game.Lock.SelectSpecialElection = false
	}

	if game.Lock.SpecialElection {
		game.Lock.SpecialElection = false
	}

	if game.Lock.SelectPartyMembershipInvestigate {
		game.Lock.SelectPartyMembershipInvestigate = false
	}

	if game.Lock.InvestigateLoyalty {
		game.Lock.InvestigateLoyalty = false
	}

	if game.Lock.ShowPlayerLoyalty {
		game.Lock.ShowPlayerLoyalty = false
	}

	if game.Lock.SelectPartyMembershipInvestigateReverse {
		game.Lock.SelectPartyMembershipInvestigateReverse = false
	}

	if game.Lock.SelectPolicies {
		game.Lock.SelectPolicies = false
	}

	if game.Lock.SelectOnePolicy {
		game.Lock.SelectOnePolicy = false
	}

	if game.Lock.SelectBurnCard {
		game.Lock.SelectBurnCard = false
	}

	game.GamePublic.GameState.PendingChancellorIndex = -1

	game.Summary.Logs = append(game.Summary.Logs, struct {
		EnactedPolicy string `bson:"enactedPolicy" json:"enactedPolicy"`
	}{
		EnactedPolicy: policy,
	})

	game.GamePublic.GeneralGameSettings.Status = "A policy is being enacted."

	if policy == "fascist" {
		game.GamePublic.TrackState.FascistPolicyCount++

	} else {
		game.GamePublic.TrackState.LiberalPolicyCount++
	}

	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	game.GamePublic.TrackState.EnactedPolicies = append(game.GamePublic.TrackState.EnactedPolicies, types.Policy{
		Position: "middle",
		CardBack: policy,
		Flipped:  false,
	})

	SendInProgressGameUpdate(game)

	timeout := 2 * time.Second

	if game.GamePublic.GeneralGameSettings.Experienced {
		timeout = 300 * time.Millisecond
	}

	time.AfterFunc(timeout, func() {
		game.GamePublic.TrackState.EnactedPolicies[len(game.GamePublic.TrackState.EnactedPolicies)-1].Flipped = true

		if policy == "fascist" {
			game.GamePublic.GameState.AudioCue = "enactPolicyF"

		} else {
			game.GamePublic.GameState.AudioCue = "enactPolicyL"
		}

		SendInProgressGameUpdate(game)
	})

	timeout = 4 * time.Second

	if game.GamePublic.GeneralGameSettings.Experienced {
		timeout = 1 * time.Second
	}

	time.AfterFunc(timeout, func() {
		game.GamePublic.GameState.AudioCue = ""

		if policy == "fascist" {
			game.GamePublic.TrackState.EnactedPolicies[len(game.GamePublic.TrackState.EnactedPolicies)-1].Position = "fascist" + strconv.Itoa(game.GamePublic.TrackState.FascistPolicyCount)

		} else {
			game.GamePublic.TrackState.EnactedPolicies[len(game.GamePublic.TrackState.EnactedPolicies)-1].Position = "liberal" + strconv.Itoa(game.GamePublic.TrackState.LiberalPolicyCount)
		}

		chat := []types.GameChat{
			types.GameChat{
				Text: "A ",
			},
			types.GameChat{
				Text: policy,
				Type: policy,
			},
		}

		if policy == "fascist" {
			chat = append(chat, types.GameChat{
				Text: " policy has been enacted. (" + strconv.Itoa(game.GamePublic.TrackState.FascistPolicyCount) + "/6)",
			})

		} else {
			chat = append(chat, types.GameChat{
				Text: " policy has been enacted. (" + strconv.Itoa(game.GamePublic.TrackState.LiberalPolicyCount) + "/5)",
			})
		}

		if !game.GamePublic.GeneralGameSettings.DisableGamechat {
			for i := range game.SeatedPlayers {
				game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.PlayerChat{
					Timestamp: time.Now(),
					GameChat:  true,
					Chat:      chat,
				})
			}

			game.UnseatedGameChats = append(game.UnseatedGameChats, types.PlayerChat{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat:      chat,
			})
		}

		if game.GamePublic.TrackState.FascistPolicyCount == 6 || game.GamePublic.TrackState.LiberalPolicyCount == 5 {
			for i := range game.GamePublic.PublicPlayerStates {
				game.GamePublic.PublicPlayerStates[i].CardStatus.CardFront = "secretrole"
				game.GamePublic.PublicPlayerStates[i].CardStatus.CardBack = game.SeatedPlayers[i].Role
				game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = true
				game.GamePublic.PublicPlayerStates[i].CardStatus.Flipped = false
			}

			SendInProgressGameUpdate(game)

			if game.GamePublic.TrackState.FascistPolicyCount == 6 {
				game.GamePublic.GameState.AudioCue = "fascistsWin"

			} else {
				game.GamePublic.GameState.AudioCue = "liberalsWin"
			}

			time.AfterFunc(2*time.Second, func() {
				for i := range game.GamePublic.PublicPlayerStates {
					game.GamePublic.PublicPlayerStates[i].CardStatus.Flipped = true
				}

				game.GamePublic.GameState.AudioCue = ""

				if game.GamePublic.TrackState.FascistPolicyCount == 6 {
					CompleteGame(game, "fascist")

				} else {
					CompleteGame(game, "liberal")
				}
			})

		} else if policy == "fascist" && game.GamePublic.CustomGameSettings.Powers[game.GamePublic.TrackState.FascistPolicyCount] != nil {
			if !game.GamePublic.GeneralGameSettings.DisableGamechat {
				gameChat := types.PlayerChat{
					Timestamp: time.Now(),
					GameChat:  true,
				}

				switch *game.GamePublic.CustomGameSettings.Powers[game.GamePublic.TrackState.FascistPolicyCount] {
				case "investigate":
					gameChat.Chat = "The president must investigate the party membership of another player."

				case "deckpeek":
					gameChat.Chat = "The president must examine the top 3 policies."

				case "election":
					gameChat.Chat = "The president must select a player for a special election."

				case "bullet":
					gameChat.Chat = "The president must select a player for execution."

				case "reverseinv":
					gameChat.Chat = "The president must reveal their party membership to another player."

				case "peekdrop":
					gameChat.Chat = "The president must examine the top policy, and may discard it."
				}

				for i := range game.SeatedPlayers {
					game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, gameChat)
				}
			}

			switch *game.GamePublic.CustomGameSettings.Powers[game.GamePublic.TrackState.FascistPolicyCount] {
			case "investigate":
				StartInvestigation(game)

			case "deckpeek":
				StartToPeekAtPolicies(game)

			case "election":
				StartSpecialElection(game)

			case "bullet":
				StartExecution(game)

			case "reverseinv":
				StartReverseInvestigation(game)

			case "peekdrop":
				StartToPeekAtPolicyToDrop(game)
			}

			for i := range game.GamePublic.PublicPlayerStates {
				if game.GamePublic.PublicPlayerStates[i].PreviousGovernmentStatus != "" {
					game.GamePublic.PublicPlayerStates[i].PreviousGovernmentStatus = ""
				}
			}

			if game.GamePublic.TrackState.ElectionTrackerCount <= 2 {
				game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].PreviousGovernmentStatus = "wasPresident"
				game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.ChancellorIndex].PreviousGovernmentStatus = "wasChancellor"
			}

			if game.GamePublic.GeneralGameSettings.Timer > 0 {
				if game.Timer != nil {
					game.Timer.Stop()
					game.Timer = nil
				}

				game.GamePublic.GameState.TimedMode = true

				game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
					if game.GamePublic.GameState.TimedMode {
						game.GamePublic.GameState.TimedMode = false
						choices := []int{}

						for i := range game.SeatedPlayers {
							if !game.SeatedPlayers[i].Dead && i != game.GamePublic.GameState.PresidentIndex {
								choices = append(choices, i)
							}
						}

						replayGameChats := []types.GameChat{
							types.GameChat{
								Text: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username,
								Type: "player",
							},
						}

						switch *game.GamePublic.CustomGameSettings.Powers[game.GamePublic.TrackState.FascistPolicyCount] {
						case "investigate":
							Investigate(game, choices[utils.RandInt(0, uint32(len(choices)))])

							replayGameChats = append(replayGameChats, types.GameChat{
								Text: " was forced by the timer to select a random player to investigate.",
							})

						case "deckpeek":
							PeekAtPolicies(game)

							replayGameChats = append(replayGameChats, types.GameChat{
								Text: " was forced by the timer to peek.",
							})

						case "election":
							SpecialElect(game, choices[utils.RandInt(0, uint32(len(choices)))])

							replayGameChats = append(replayGameChats, types.GameChat{
								Text: " was forced by the timer to select a random player to special elect.",
							})

						case "bullet":
							if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].Role.CardName == "fascist" {
								for i := range choices {
									if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].Role.CardName == "hitler" {
										choices = append([]int{}, append(choices[:i], choices[i+1:]...)...)
									}
								}
							}

							Execute(game, choices[utils.RandInt(0, uint32(len(choices)))])

							replayGameChats = append(replayGameChats, types.GameChat{
								Text: " was forced by the timer to select a random player to execute.",
							})

							// case "reverseinv":
							// 	ReverseInvestigate(game, choices[utils.RandInt(0, uint32(len(choices)))])

							// case "peekdrop":
							// 	PeekAtPolicyToDrop(game)
						}

						game.ReplayGameChats = append(game.ReplayGameChats, types.PlayerChat{
							Timestamp: time.Now(),
							GameChat:  true,
							Chat:      replayGameChats,
						})
					}
				})

				SendInProgressGameUpdate(game)
			}

		} else {
			SendInProgressGameUpdate(game)

			for i := range game.GamePublic.PublicPlayerStates {
				if game.GamePublic.PublicPlayerStates[i].PreviousGovernmentStatus != "" {
					game.GamePublic.PublicPlayerStates[i].PreviousGovernmentStatus = ""
				}
			}

			if game.GamePublic.TrackState.ElectionTrackerCount <= 2 {
				game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].PreviousGovernmentStatus = "wasPresident"
				game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.ChancellorIndex].PreviousGovernmentStatus = "wasChancellor"
			}

			StartElection(game, -1)
		}

		game.GamePublic.TrackState.ElectionTrackerCount = 0
	})
}

func SelectPresidentPolicy(user *types.UserPublic, game *types.GamePrivate, index int, timer bool) {
	if user == nil || game == nil {
		return
	}

	if game.GamePublic.GameState.Frozen {
		fmt.Println("Frozen")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+user.ID, "sendAlert", "An AEM member has prevented this game from proceeding. Please wait.")
		return
	}

	if game.GamePublic.GeneralGameSettings.Remade {
		fmt.Println("Remade")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+user.ID, "sendAlert", "This game has been remade and is now no longer playable.")
		return
	}

	if user.ID != game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.ID {
		return
	}

	fmt.Println("Selecting policy", index)

	game.ActionMutex.Lock()

	if game.Lock.SelectPresidentPolicy {
		fmt.Println("President already selected policy")
		game.ActionMutex.Unlock()
		return
	}

	game.Lock.SelectPresidentPolicy = true

	hiddenInfoChat := []types.GameChat{
		types.GameChat{
			Text: "President ",
		},
		types.GameChat{
			Text: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PresidentIndex+1) + "}",
			Type: "player",
		},
	}

	if timer {
		hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
			Text: " has automatically discarded a ",
		},
			types.GameChat{
				Text: game.CurrentElectionPolicies[index],
				Type: game.CurrentElectionPolicies[index],
			},
			types.GameChat{
				Text: " policy due to the timer expiring.",
			})

	} else {
		hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
			Text: " has selected a ",
		},
			types.GameChat{
				Text: game.CurrentElectionPolicies[index],
				Type: game.CurrentElectionPolicies[index],
			},
			types.GameChat{
				Text: " policy.",
			})
	}

	modChat := types.PlayerChat{
		Timestamp: time.Now(),
		GameChat:  true,
		Chat:      hiddenInfoChat,
	}

	game.HiddenInfoChat = append(game.HiddenInfoChat, modChat)
	SendInProgressModChatUpdate(game, modChat)

	game.CurrentChancellorOptions = make([]string, index)
	copy(game.CurrentChancellorOptions, game.CurrentElectionPolicies[:index])
	game.CurrentChancellorOptions = append(game.CurrentChancellorOptions, game.CurrentElectionPolicies[index+1:]...)

	if !timer && !game.GamePublic.GeneralGameSettings.Private {
		gameType := "Ranked"

		if game.GamePublic.GeneralGameSettings.Casual {
			gameType = "Casual"

		} else if game.GamePublic.GeneralGameSettings.Practice {
			gameType = "Practice"
		}

		comment := ""

		if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].Role.CardName == "fascist" && game.CurrentElectionPolicies[index] == "fascist" {
			if game.GamePublic.TrackState.LiberalPolicyCount == 4 {
				if game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].Role.Team == "fascist" && game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "liberal" {
					comment = "got BBR with 4 blues on the track and forced blues on a fascist chancellor."

				} else if game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].Role.Team != "fascist" && ((game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "liberal") || (game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "fascist")) {
					comment = "got BRR with 4 blues on the track and offered choice to a liberal chancellor."
				}

			} else if game.GamePublic.TrackState.FascistPolicyCount == 5 {
				if game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].Role.Team == "fascist" && game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "liberal" {
					comment = "got BBR with 5 reds on the track and forced blues on a fascist chancellor."

				} else if game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].Role.Team != "fascist" && ((game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "liberal") || (game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "fascist")) {
					comment = "got BRR with 5 reds on the track and offered choice to a liberal chancellor."
				}
			}

		} else if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].Role.CardName == "hitler" && game.CurrentElectionPolicies[index] == "fascist" {
			if game.GamePublic.TrackState.LiberalPolicyCount == 4 {
				if game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "liberal" {
					comment = "got BBR with 4 blues on the track and forced blues."

				} else if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[game.GamePublic.GameState.ChancellorIndex].Role.Team != "fascist" && ((game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "liberal") || (game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "fascist")) {
					if game.GamePublic.TrackState.FascistPolicyCount == 5 {
						comment = "got BRR with 4 blues on the track and did not force 6th red."

					} else {
						comment = "got BRR with 4 blues on the track and offered choice."
					}
				}

			} else if game.GamePublic.TrackState.FascistPolicyCount == 5 && game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[game.GamePublic.GameState.ChancellorIndex].Role.Team != "fascist" && ((game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "liberal") || (game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "fascist")) {
				comment = "got BRR with 5 reds on the track and did not force 6th red."
			}

		} else if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].Role.CardName == "liberal" && game.CurrentElectionPolicies[index] == "liberal" {
			if game.GamePublic.TrackState.FascistPolicyCount == 5 {
				if game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "fascist" {
					comment = "got BRR during veto zone and forced reds."

				} else if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[game.GamePublic.GameState.ChancellorIndex].Role.Team != "liberal" && ((game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "liberal") || (game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "fascist")) { // && game.GamePublic.TrackState.LiberalPolicyCount == 4?
					if game.GamePublic.TrackState.LiberalPolicyCount == 4 {
						comment = "got BBR during veto zone and did not force 5th blue."

					} else {
						comment = "got BBR during veto zone and offered choice."
					}
				}

			} else if game.GamePublic.TrackState.LiberalPolicyCount == 4 {
				if game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "fascist" {
					comment = "got BRR with 4 blues on the track and forced reds."

				} else if game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[game.GamePublic.GameState.ChancellorIndex].Role.Team != "liberal" && ((game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "liberal") || (game.CurrentChancellorOptions[0] == "liberal" && game.CurrentChancellorOptions[1] == "fascist")) {
					comment = "got BBR with 4 blues on the track and did not force 5th blue."
				}

			} else if game.GamePublic.TrackState.FascistPolicyCount < 3 {
				if game.CurrentChancellorOptions[0] == "fascist" && game.CurrentChancellorOptions[1] == "fascist" {
					comment = "got BRR before HZ and forced reds."
				}
			}
		}

		if comment != "" {
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
	}

	game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].Loader = false
	game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.ChancellorIndex].Loader = true

	newCardFlingerState := make([]types.CardFlinger, index)
	copy(newCardFlingerState, game.CardFlingerState[:index])
	newCardFlingerState = append(newCardFlingerState, game.CardFlingerState[index+1:]...)
	newCardFlingerState[0].Position = "middle-left"
	newCardFlingerState[1].Position = "middle-right"

	for i := range newCardFlingerState {
		newCardFlingerState[i].CardStatus.Flipped = false
		newCardFlingerState[i].NotificationStatus = ""
	}

	for i := range game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState {
		if i == index {
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState[i].NotificationStatus = "selected"

		} else {
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState[i].NotificationStatus = ""
		}

		game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState[i].Action = ""
		game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState[i].CardStatus.Flipped = false
	}

	game.Summary.Logs = append(game.Summary.Logs, struct {
		ChancellorHand []string `bson:"chancellorHand" json:"chancellorHand"`
	}{
		ChancellorHand: game.CurrentChancellorOptions,
	})

	game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState = newCardFlingerState
	game.GamePublic.GeneralGameSettings.Status = "Waiting on chancellor enactment."
	game.GamePublic.GameState.Phase = "chancellorSelectingPolicy"

	fmt.Println("Experienced, DisableGamechat", game.GamePublic.GeneralGameSettings.Experienced, game.GamePublic.GeneralGameSettings.DisableGamechat)

	if !game.GamePublic.GeneralGameSettings.Experienced && !game.GamePublic.GeneralGameSettings.DisableGamechat {
		game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].GameChats = append(game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].GameChats, types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: "As chancellor, you must select a policy to enact.",
				},
			},
		})
	}

	SendInProgressGameUpdate(game)

	time.AfterFunc(200*time.Millisecond, func() {
		game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].CardFlingerState = []types.CardFlinger{}

		for i := range game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState {
			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].CardStatus.Flipped = true
			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].NotificationStatus = "notification"
		}

		if game.GamePublic.GeneralGameSettings.Timer > 0 {
			if game.Timer != nil {
				game.Timer.Stop()
				game.Timer = nil
			}

			game.GamePublic.GameState.TimedMode = true

			game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
				if game.GamePublic.GameState.TimedMode {
					SelectChancellorPolicy(&game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].UserPublic, game, int(utils.RandInt(0, uint32(len(game.CurrentChancellorOptions)))), true)
				}
			})
		}

		fmt.Println("200ms passed")

		SendInProgressGameUpdate(game)
	})

	game.ActionMutex.Unlock()
}

func SelectChancellorPolicy(user *types.UserPublic, game *types.GamePrivate, index int, timer bool) {
	if user == nil || game == nil {
		return
	}

	if game.GamePublic.GameState.Frozen {
		fmt.Println("Frozen")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+user.ID, "sendAlert", "An AEM member has prevented this game from proceeding. Please wait.")
		return
	}

	if game.GamePublic.GeneralGameSettings.Remade {
		fmt.Println("Remade")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+user.ID, "sendAlert", "This game has been remade and is now no longer playable.")
		return
	}

	if user.ID != game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].UserPublic.ID {
		return
	}

	fmt.Println("Chancellor selected", index)

	game.ActionMutex.Lock()

	if game.Lock.SelectChancellorPolicy || len(game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState) == 0 {
		fmt.Println("Chancellor already selected policy")
		game.ActionMutex.Unlock()
		return
	}

	game.Lock.SelectPresidentPolicy = false
	game.Lock.SelectChancellorPolicy = true
	fmt.Println("*Chancellor enacted", index)

	if !timer && !game.GamePublic.GeneralGameSettings.Private {
		gameType := "Ranked"

		if game.GamePublic.GeneralGameSettings.Casual {
			gameType = "Casual"

		} else if game.GamePublic.GeneralGameSettings.Practice {
			gameType = "Practice"
		}

		comment := ""

		if game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].Role.CardName == "fascist" && game.GamePublic.TrackState.LiberalPolicyCount == 4 && game.CurrentChancellorOptions[index] == "liberal" && (game.CurrentChancellorOptions[0] == "fascist" || game.CurrentChancellorOptions[1] == "fascist") {
			comment = "was given choice as chancellor with 4 blues on the track, and played liberal."

		} else if game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].Role.CardName == "liberal" && game.CurrentChancellorOptions[index] == "fascist" && (game.CurrentChancellorOptions[0] == "liberal" || game.CurrentChancellorOptions[1] == "liberal") {
			comment = "was given choice as chancellor, and played fascist."
		}

		if comment != "" {
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
	}

	hiddenInfoChat := []types.GameChat{
		types.GameChat{
			Text: "Chancellor ",
		},
		types.GameChat{
			Text: game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.ChancellorIndex+1) + "}",
			Type: "player",
		},
	}

	if timer {
		hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
			Text: " has automatically chosen to play a ",
		},
			types.GameChat{
				Text: game.CurrentChancellorOptions[index],
				Type: game.CurrentChancellorOptions[index],
			},
			types.GameChat{
				Text: " policy due to the timer expiring.",
			})

	} else {
		hiddenInfoChat = append(hiddenInfoChat, types.GameChat{
			Text: " has chosen to play a ",
		},
			types.GameChat{
				Text: game.CurrentChancellorOptions[index],
				Type: game.CurrentChancellorOptions[index],
			},
			types.GameChat{
				Text: " policy.",
			})
	}

	modChat := types.PlayerChat{
		Timestamp: time.Now(),
		GameChat:  true,
		Chat:      hiddenInfoChat,
	}

	game.HiddenInfoChat = append(game.HiddenInfoChat, modChat)
	SendInProgressModChatUpdate(game, modChat)

	/*if game.GamePublic.GeneralGameSettings.Timer > 0 && game.Timer != nil {
	 	game.Timer.Stop()
	 	game.Timer = nil
	 	game.GameState.TimedMode = false
	}*/

	for i := range game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState {
		if i == index {
			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].NotificationStatus = "Selected"

		} else {
			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].NotificationStatus = ""
		}

		game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].Action = ""
		game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].CardStatus.Flipped = false
	}

	game.GamePublic.PublicPlayerStates[game.GameState.ChancellorIndex].Loader = false

	if game.GamePublic.GameState.Veto {
		game.CurrentElectionPolicies = make([]string, index)
		copy(game.CurrentElectionPolicies, game.CurrentChancellorOptions[:index])
		game.CurrentElectionPolicies = append(game.CurrentElectionPolicies, game.CurrentChancellorOptions[index+1:]...)
		game.GamePublic.GeneralGameSettings.Status = "Chancellor to vote on policy veto."
		SendInProgressGameUpdate(game)

		time.AfterFunc(1*time.Second, func() {
			game.GamePublic.PublicPlayerStates[game.GameState.ChancellorIndex].Loader = true

			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState = []types.CardFlinger{
				types.CardFlinger{
					Position:           "middle-left",
					NotificationStatus: "",
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   false,
						CardFront: "ballot",
						CardBack:  "ja",
					},
				},
				types.CardFlinger{
					Position:           "middle-right",
					NotificationStatus: "",
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   false,
						CardFront: "ballot",
						CardBack:  "nein",
					},
				},
			}
		})

		game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].GameChats = append(game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].GameChats, types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: "You must vote whether or not to veto these policies.  Select Ja to veto the your chosen policy or select Nein to enact your chosen policy.",
				},
			},
		})

		SendInProgressGameUpdate(game)

		time.AfterFunc(500*time.Millisecond, func() {
			for i := range game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState {
				game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].CardStatus.Flipped = true
				game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState[i].NotificationStatus = "notification"
				game.GamePublic.GameState.Phase = "chancellorVoteOnVeto"

				if game.GamePublic.GeneralGameSettings.Timer > 0 {
					if game.Timer != nil {
						game.Timer.Stop()
						game.Timer = nil
					}

					game.GamePublic.GameState.TimedMode = true

					game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
						if game.GamePublic.GameState.TimedMode {
							game.GamePublic.GameState.TimedMode = false
							SelectChancellorVoteOnVeto(&game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].UserPublic, game, utils.RandInt(0, 1) == 1)

							game.ReplayGameChats = append(game.ReplayGameChats, types.PlayerChat{
								Timestamp: time.Now(),
								GameChat:  true,
								Chat: []types.GameChat{
									types.GameChat{
										Text: game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].UserPublic.Username,
										Type: "player",
									},
									types.GameChat{
										Text: " was forced by the timer to select a random veto vote.",
									},
								},
							})
						}
					})
				}
			}
		})

	} else {
		game.CurrentElectionPolicies = []string{}
		game.GamePublic.GameState.Phase = "enactPolicy"
		SendInProgressGameUpdate(game)

		time.AfterFunc(200*time.Millisecond, func() {
			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].CardFlingerState = []types.CardFlinger{}
			EnactPolicy(game, game.CurrentChancellorOptions[index])
		})
	}

	if game.GamePublic.GeneralGameSettings.Experienced {
		game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[game.GamePublic.GameState.PresidentIndex].Claim = "wasPresident"
		game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].PlayerStates[game.GamePublic.GameState.ChancellorIndex].Claim = "wasChancellor"

	} else {
		time.AfterFunc(3*time.Second, func() {
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[game.GamePublic.GameState.PresidentIndex].Claim = "wasPresident"
			game.SeatedPlayers[game.GamePublic.GameState.ChancellorIndex].PlayerStates[game.GamePublic.GameState.ChancellorIndex].Claim = "wasChancellor"
			SendInProgressGameUpdate(game)
		})
	}

	game.ActionMutex.Unlock()
}
