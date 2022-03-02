package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
	// "github.com/googollee/go-socket.io"
)

func SelectChancellor(user *types.UserPublic, game *types.GamePrivate, chancellorIndex int, force bool) {
	// user := GetUser(socket)

	if user == nil || game == nil {
		return
	}

	fmt.Println("SelectChancellor", user.ID, game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.ID, chancellorIndex, force)

	if user.ID != game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.ID {
		return
	}

	validChoice := false
	clickActionInfo, ok := game.GamePublic.GameState.ClickActionInfo[1].([]int)

	if !ok {
		return
	}

	for i := range clickActionInfo {
		// choice, ok := game.GamePublic.GameState.ClickActionInfo[1][i].(int)
		fmt.Println("OK, CI, C", ok, chancellorIndex, clickActionInfo[i])

		if ok && chancellorIndex == clickActionInfo[i] {
			validChoice = true
		}
	}

	if (game.GeneralGameSettings.Tourny && game.GeneralGameSettings.TournyInfo.Cancelled) || !validChoice {
		return
	}

	fmt.Println("OK")

	if game.GamePublic.GameState.Frozen && !force {
		fmt.Println("Frozen")
		IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-"+user.ID, "sendAlert", "An AEM member has prevented this game from proceeding. Please wait.")
		// socket.Emit("sendAlert", "An AEM member has prevented this game from proceeding. Please wait.")
		return
	}

	if game.GamePublic.GameState.TimedMode && game.Timer != nil {
		game.Timer.Stop()
		game.Timer = nil
		game.GamePublic.GameState.TimedMode = false
	}

	game.ActionMutex.Lock()
	fmt.Println("Taking action", game.GamePublic.GameState.PendingChancellorIndex, game.GamePublic.GameState.Phase)

	if game.GamePublic.GameState.PendingChancellorIndex == -1 && game.GamePublic.GameState.Phase != "voting" {
		game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].Loader = false

		game.Summary.Logs = append(game.Summary.Logs, struct {
			ChancellorIndex int `bson:"chancellorIndex" json:"chancellorId"`
		}{
			ChancellorIndex: chancellorIndex,
		})

		for i := range game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates {
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[i].NotificationStatus = ""
		}

		game.GamePublic.PublicPlayerStates[chancellorIndex].GovernmentStatus = "isPendingChancellor"
		game.GamePublic.GameState.PendingChancellorIndex = chancellorIndex
		game.GeneralGameSettings.Status = "Vote on election #" + strconv.Itoa(game.GeneralGameSettings.ElectionCount) + " now."

		for i := range game.GamePublic.PublicPlayerStates {
			if !game.GamePublic.PublicPlayerStates[i].Dead {
				game.GamePublic.PublicPlayerStates[i].Loader = true

				game.GamePublic.PublicPlayerStates[i].CardStatus = types.CardStatus{
					CardDisplayed: true,
					Flipped:       false,
					CardFront:     "ballot",
					CardBack:      "",
				}
			}
		}

		SendInProgressGameUpdate(game)

		for i := range game.SeatedPlayers {
			if !game.GamePublic.GeneralGameSettings.DisableGamechat {
				game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.PlayerChat{
					Timestamp: time.Now(),
					GameChat:  true,
					Chat: []types.GameChat{
						types.GameChat{
							Text: "You must vote for the election of president ",
						},
						types.GameChat{
							Text: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PresidentIndex+1) + "}",
							Type: "player",
						},
						types.GameChat{
							Text: " and chancellor ",
						},
						types.GameChat{
							Text: game.SeatedPlayers[game.GamePublic.GameState.PendingChancellorIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PendingChancellorIndex+1) + "}",
							Type: "player",
						},
						types.GameChat{
							Text: ".",
						},
					},
				})
			}

			game.SeatedPlayers[i].CardFlingerState = []types.CardFlingerState{
				types.CardFlingerState{
					Position:           "middle-left",
					NotificationStatus: "",
					Action:             "active",
					CardStatus: types.CardStatus{
						Flipped:   false,
						CardFront: "ballot",
						CardBack:  "ja",
					},
				},
				types.CardFlingerState{
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

			fmt.Println("Card Flinger State", game.SeatedPlayers[i].CardFlingerState)
		}

		game.UnseatedGameChats = append(game.UnseatedGameChats, types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: "President ",
				},
				types.GameChat{
					Text: game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PresidentIndex+1) + "}",
					Type: "player",
				},
				types.GameChat{
					Text: " nominates ",
				},
				types.GameChat{
					Text: game.SeatedPlayers[game.GamePublic.GameState.PendingChancellorIndex].UserPublic.Username + " {" + strconv.Itoa(game.GamePublic.GameState.PendingChancellorIndex+1) + "}",
					Type: "player",
				},
				types.GameChat{
					Text: " as chancellor.",
				},
			},
		})

		time.AfterFunc(500*time.Millisecond, func() {
			SendInProgressGameUpdate(game)
		})

		game.GamePublic.GameState.Phase = "voting"

		time.AfterFunc(500*time.Millisecond, func() {
			for i := range game.SeatedPlayers {
				if len(game.SeatedPlayers[i].CardFlingerState) > 0 {
					game.SeatedPlayers[i].CardFlingerState[0].CardStatus.Flipped = true
					game.SeatedPlayers[i].CardFlingerState[1].CardStatus.Flipped = true
					game.SeatedPlayers[i].VoteStatus.Voted = false
				}
			}

			if game.GamePublic.GeneralGameSettings.Timer > 0 {
				if game.Timer != nil {
					game.Timer.Stop()
					game.Timer = nil
				}

				game.GamePublic.GameState.TimedMode = true

				game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
					fmt.Println("Vote Timer Expired")

					neededPlayers := game.GamePublic.PlayerCount/2 + 2
					activePlayerCount := 0

					for i := range game.GamePublic.PublicPlayerStates {
						if !game.GamePublic.PublicPlayerStates[i].LeftGame || game.GamePublic.PublicPlayerStates[i].Dead {
							activePlayerCount++
						}
					}

					if activePlayerCount < neededPlayers {
						if !game.GamePublic.GeneralGameSettings.DisableGamechat {
							for i := range game.SeatedPlayers {
								game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.PlayerChat{
									Timestamp: time.Now(),
									GameChat:  true,
									Chat: []types.GameChat{
										types.GameChat{
											Text: "Not enough players are present, votes will not be auto-picked.",
										},
									},
								})
							}

							SendInProgressGameUpdate(game)
						}

						return
					}

					if game.GamePublic.GameState.TimedMode {
						game.GamePublic.GameState.TimedMode = false

						for i := range game.GamePublic.PublicPlayerStates {
							if !game.SeatedPlayers[i].VoteStatus.Voted && !game.SeatedPlayers[i].Dead {
								SelectVote(&game.SeatedPlayers[i], game, utils.RandInt(0, 1) == 1, false)
							}
						}
					}
				})
			}
		})

		SendInProgressGameUpdate(game)
	}
	game.ActionMutex.Unlock()
}
