package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
)

func BeginGame(game *types.GamePrivate) {
	fmt.Println("Game is starting...")

	game.GamePublic.GeneralGameSettings.TimeStarted = time.Now()
	customGameSettings := game.GamePublic.CustomGameSettings

	if !game.GamePublic.CustomGameSettings.Enabled {
		game.GamePublic.CustomGameSettings.HitlerZone = 3
		game.GamePublic.CustomGameSettings.VetoZone = 5
		game.GamePublic.CustomGameSettings.TrackState = types.CustomGameCounter{
			Liberal: 0,
			Fascist: 0,
		}
		game.GamePublic.CustomGameSettings.DeckState = types.CustomGameCounter{
			Liberal: 6,
			Fascist: 11,
		}

		if 5 <= game.GamePublic.PlayerCount && game.GamePublic.PlayerCount <= 6 {
			game.CustomGameSettings.FascistCount = 1
			game.CustomGameSettings.HitlerKnowsFascists = true
			game.CustomGameSettings.Powers = []string{"null", "null", "deckpeek", "bullet", "bullet"}

		} else if 7 <= game.GamePublic.PlayerCount && game.GamePublic.PlayerCount <= 8 {
			game.CustomGameSettings.FascistCount = 2
			game.CustomGameSettings.HitlerKnowsFascists = false
			game.CustomGameSettings.Powers = []string{"null", "investigate", "election", "bullet", "bullet"}

		} else if 9 <= game.GamePublic.PlayerCount && game.GamePublic.PlayerCount <= 10 {
			game.CustomGameSettings.FascistCount = 3
			game.CustomGameSettings.HitlerKnowsFascists = false
			game.CustomGameSettings.Powers = []string{"investigate", "investigate", "election", "bullet", "bullet"}
		}

		game.GamePublic.CustomGameSettings.LiberalCount = game.GamePublic.PlayerCount - customGameSettings.FascistCount - 1
	}

	ShufflePolicies(game)

	fmt.Println("Custom Game Settings", customGameSettings, game.GamePublic.CustomGameSettings)

	roles := make([]types.CardBack, game.GamePublic.PlayerCount)

	roles[0] = types.CardBack{
		CardName: "hitler",
		Icon:     1,
		Team:     "fascist",
	}

	for i := 0; i < game.GamePublic.CustomGameSettings.LiberalCount; i++ {
		roles[i+1] = types.CardBack{
			CardName: "liberal",
			Icon:     i % 6,
			Team:     "liberal",
		}
	}

	for i := 6; i < 6+game.GamePublic.CustomGameSettings.FascistCount; i++ {
		roles[i+game.GamePublic.CustomGameSettings.LiberalCount-6] = types.CardBack{
			CardName: "fascist",
			Icon:     i,
			Team:     "fascist",
		}
	}

	// fmt.Println("Roles", roles)

	for i := 0; i < game.GamePublic.PlayerCount-1; i++ {
		j := utils.RandInt(uint32(i), uint32(game.GamePublic.PlayerCount))
		role := roles[i]
		roles[i] = roles[j]
		roles[j] = role
	}

	// fmt.Println("Shuffled Roles", roles)

	game.GamePublic.GeneralGameSettings.Status = "Dealing roles..."

	for i := range game.GamePublic.PublicPlayerStates {
		game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = true
		game.GamePublic.PublicPlayerStates[i].CardStatus.CardFront = "secretrole"
	}

	// fmt.Println("Player States:", game.GamePublic.PublicPlayerStates, game.GamePublic.PublicPlayerStates[0].CardStatus)

	liberalPlayers := make([]types.PlayerState, game.GamePublic.CustomGameSettings.LiberalCount)
	liberalPlayerCount := 0
	var hitler types.PlayerState
	fascistPlayers := make([]types.PlayerState, game.GamePublic.CustomGameSettings.FascistCount)
	fascistPlayerCount := 0

	for i := range game.SeatedPlayers {
		game.SocketMapMutex.RLock()
		for j := range game.SocketMap[game.SeatedPlayers[i].ID] {
			IO.LeaveRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID+"-observer", game.SocketMap[game.SeatedPlayers[i].ID][j])
		}
		game.SocketMapMutex.RUnlock()

		game.SeatedPlayers[i].Role = roles[i]
		game.SeatedPlayers[i].PlayerStates = make([]types.PlayerState, game.GamePublic.PlayerCount)

		for j := 0; j < game.GamePublic.PlayerCount; j++ {
			game.SeatedPlayers[i].PlayerStates[j] = types.PlayerState{
				CardStatus: types.CardStatus{
					CardBack:      "",
					CardDisplayed: true,
					CardFront:     "secretrole",
				},
				NotificationStatus: "",
				NameStatus:         "",
			}
		}

		game.SeatedPlayers[i].PlayerStates[i].CardStatus.CardBack = game.SeatedPlayers[i].Role
		game.SeatedPlayers[i].PlayerStates[i].CardStatus.CardName = game.SeatedPlayers[i].Role.CardName
		// game.SeatedPlayers[i].PlayerStates[i].CardStatus.CardDisplayed = true
		game.SeatedPlayers[i].PlayerStates[i].CardStatus.Flipped = true
		game.SeatedPlayers[i].PlayerStates[i].NameStatus = game.SeatedPlayers[i].Role.CardName
		game.SeatedPlayers[i].PlayerStates[i].NotificationStatus = game.SeatedPlayers[i].Role.CardName
		// fmt.Println("Player States:", i, game.SeatedPlayers[i].PlayerStates[i].CardStatus)

		if !game.GamePublic.GeneralGameSettings.DisableGamechat {
			game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.PlayerChat{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat: []types.GameChat{
					types.GameChat{
						Text: "The game begins and you receive the ",
					},
					types.GameChat{
						Text: game.SeatedPlayers[i].Role.CardName,
						Type: game.SeatedPlayers[i].Role.CardName,
					},
					types.GameChat{
						Text: " role and take seat #",
					},
					types.GameChat{
						Text: strconv.Itoa(i+1) + ".",
						Type: "player",
					},
				},
			})

		} else {
			game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, types.PlayerChat{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat: []types.GameChat{
					types.GameChat{
						Text: "The game begins.",
					},
				},
			})
		}

		game.HiddenInfoChat = append(game.HiddenInfoChat, types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: game.SeatedPlayers[i].UserPublic.Username + " {" + strconv.Itoa(i+1) + "}",
					Type: "player",
				},
				types.GameChat{
					Text: " is assigned the ",
				},
				types.GameChat{
					Text: game.SeatedPlayers[i].Role.CardName,
					Type: game.SeatedPlayers[i].Role.CardName,
				},
				types.GameChat{
					Text: " role.",
				},
			},
		})

		if game.SeatedPlayers[i].Role.CardName == "fascist" {
			fascistPlayers[fascistPlayerCount] = game.SeatedPlayers[i]
			fascistPlayerCount++

		} else if game.SeatedPlayers[i].Role.CardName == "hitler" {
			hitler = game.SeatedPlayers[i]

		} else {
			liberalPlayers[liberalPlayerCount] = game.SeatedPlayers[i]
			liberalPlayerCount++
		}
	}

	fascistElo := types.TeamElo{
		Overall:  1600,
		Seasonal: 1600,
	}

	liberalElo := types.TeamElo{
		Overall:  1600,
		Seasonal: 1600,
	}

	for i := range fascistPlayers {
		fascistElo.Overall += fascistPlayers[i].UserPublic.EloOverall
		fascistElo.Seasonal += fascistPlayers[i].UserPublic.EloSeason
	}

	fascistElo.Overall += hitler.UserPublic.EloOverall
	fascistElo.Seasonal += hitler.UserPublic.EloSeason

	for i := range liberalPlayers {
		liberalElo.Overall += liberalPlayers[i].UserPublic.EloOverall
		liberalElo.Seasonal += liberalPlayers[i].UserPublic.EloSeason
	}

	fascistElo.Overall /= float64(fascistPlayerCount + 1)
	fascistElo.Seasonal /= float64(fascistPlayerCount + 1)
	liberalElo.Overall /= float64(liberalPlayerCount)
	liberalElo.Seasonal /= float64(liberalPlayerCount)

	game.Summary = types.GameSummary{
		GameID:             game.GamePublic.GeneralGameSettings.ID,
		Time:               time.Now(),
		GameSettings:       game.GamePublic.GeneralGameSettings,
		CustomGameSettings: game.GamePublic.CustomGameSettings,
		Logs:               []interface{}{},
		LiberalElo:         liberalElo,
		FascistElo:         fascistElo,
	}

	for i := range game.SeatedPlayers {
		game.Summary.Logs = append(game.Summary.Logs, struct {
			Username string `bson:"username" json:"username"`
			Role     string `bson:"role" json:"role"`
			Icon     int    `bson:"icon" json:"icon"`
		}{
			Username: game.SeatedPlayers[i].UserPublic.Username,
			Role:     game.SeatedPlayers[i].Role.CardName,
			Icon:     game.SeatedPlayers[i].Role.Icon,
		})
	}

	game.UnseatedGameChats = append(game.UnseatedGameChats, types.PlayerChat{
		Timestamp: time.Now(),
		GameChat:  true,
		Chat: []types.GameChat{
			types.GameChat{
				Text: "The game begins.",
			},
		},
	})

	// fmt.Println("Seated Players (Start Game)", game.SeatedPlayers)
	SendInProgressGameUpdate(game)

	time.AfterFunc(2000*time.Millisecond, func() {
		for i := range fascistPlayers {
			if !game.GamePublic.GeneralGameSettings.DisableGamechat {
				if game.GamePublic.CustomGameSettings.FascistCount >= 2 {
					newGameChat := []types.GameChat{
						types.GameChat{
							Text: "You see that the other ",
						},
					}

					if game.GamePublic.CustomGameSettings.FascistCount == 2 {
						newGameChat = append(newGameChat, types.GameChat{
							Text: "fascist",
							Type: "fascist",
						}, types.GameChat{
							Text: " in this game is ",
						})

					} else {
						newGameChat = append(newGameChat, types.GameChat{
							Text: "fascists",
							Type: "fascist",
						}, types.GameChat{
							Text: " in this game are ",
						})
					}

					for j := range fascistPlayers {
						if j != i {
							newGameChat = append(newGameChat, types.GameChat{
								Text: fascistPlayers[j].UserPublic.Username + " {" + strconv.Itoa(fascistPlayers[j].Index+1) + "}",
								Type: "player",
							})

							if j == len(fascistPlayers)-2 ||
								(j == len(fascistPlayers)-3 && i == len(fascistPlayers)-1) {
								/* ||
								(j == len(fascistPlayers)-4 &&
									((i == len(fascistPlayers)-2 && fascistPlayers[len(fascistPlayers)-1].NameStatus == "hitler") ||
										(i == len(fascistPlayers)-1 && fascistPlayers[len(fascistPlayers)-2].NameStatus == "hitler"))) {*/
								newGameChat = append(newGameChat, types.GameChat{
									Text: " and ",
								})
							}
						}
					}

					game.SeatedPlayers[fascistPlayers[i].Index].GameChats = append(game.SeatedPlayers[fascistPlayers[i].Index].GameChats, types.PlayerChat{
						Timestamp: time.Now(),
						GameChat:  true,
						Chat:      newGameChat,
					})
				}

				newGameChat := []types.GameChat{
					types.GameChat{
						Text: "You see that ",
					},
					types.GameChat{
						Text: "hitler",
						Type: "hitler",
					},
					types.GameChat{
						Text: " in this game is ",
					},
					types.GameChat{
						Text: hitler.UserPublic.Username + " {" + strconv.Itoa(hitler.Index+1) + "}",
						Type: "player",
					},
				}

				if game.GamePublic.CustomGameSettings.HitlerKnowsFascists {
					newGameChat = append(newGameChat, types.GameChat{
						Text: ". They also see that you are a ",
					})

				} else {
					newGameChat = append(newGameChat, types.GameChat{
						Text: ". They do not know that you are a ",
					})
				}

				newGameChat = append(newGameChat, []types.GameChat{
					types.GameChat{
						Text: "fascist",
						Type: "fascist",
					},
					types.GameChat{
						Text: ".",
					},
				}...)

				game.SeatedPlayers[fascistPlayers[i].Index].GameChats = append(game.SeatedPlayers[fascistPlayers[i].Index].GameChats, types.PlayerChat{
					Timestamp: time.Now(),
					GameChat:  true,
					Chat:      newGameChat,
				})
			}

			for j := range fascistPlayers {
				game.SeatedPlayers[fascistPlayers[i].Index].PlayerStates[fascistPlayers[j].Index].NameStatus = "fascist"
				game.SeatedPlayers[fascistPlayers[i].Index].PlayerStates[fascistPlayers[j].Index].NotificationStatus = "fascist"
			}

			game.SeatedPlayers[fascistPlayers[i].Index].PlayerStates[hitler.Index].NameStatus = "hitler"
			game.SeatedPlayers[fascistPlayers[i].Index].PlayerStates[hitler.Index].NotificationStatus = "hitler"
		}

		if !game.GamePublic.GeneralGameSettings.DisableGamechat {
			var newGameChat []types.GameChat

			if game.GamePublic.CustomGameSettings.HitlerKnowsFascists {
				newGameChat = []types.GameChat{
					types.GameChat{
						Text: "You see that the other ",
					},
				}

				if game.GamePublic.CustomGameSettings.FascistCount == 1 {
					newGameChat = append(newGameChat, types.GameChat{
						Text: "fascist",
						Type: "fascist",
					}, types.GameChat{
						Text: " in this game is ",
					})

				} else {
					newGameChat = append(newGameChat, types.GameChat{
						Text: "fascists",
						Type: "fascist",
					}, types.GameChat{
						Text: " in this game are ",
					})
				}

				for i := range fascistPlayers {
					newGameChat = append(newGameChat, types.GameChat{
						Text: fascistPlayers[i].UserPublic.Username + " {" + strconv.Itoa(fascistPlayers[i].Index+1) + "}",
						Type: "player",
					})

					if i == len(fascistPlayers)-2 {
						newGameChat = append(newGameChat, types.GameChat{
							Text: " and ",
						})
					}
				}

			} else {
				if game.GamePublic.CustomGameSettings.FascistCount == 1 {
					newGameChat = []types.GameChat{
						types.GameChat{
							Text: "There is ",
						},
						types.GameChat{
							Text: "1 fascist",
							Type: "fascist",
						},
					}

				} else {
					newGameChat = []types.GameChat{
						types.GameChat{
							Text: "There are ",
						},
						types.GameChat{
							Text: strconv.Itoa(game.GamePublic.CustomGameSettings.FascistCount) + " fascists",
							Type: "fascist",
						},
					}
				}

				newGameChat = append(newGameChat, types.GameChat{
					Text: ". They know who you are.",
				})
			}

			game.SeatedPlayers[hitler.Index].GameChats = append(game.SeatedPlayers[hitler.Index].GameChats, types.PlayerChat{
				Timestamp: time.Now(),
				GameChat:  true,
				Chat:      newGameChat,
			})
		}

		if game.GamePublic.CustomGameSettings.HitlerKnowsFascists {
			for i := range fascistPlayers {
				game.SeatedPlayers[hitler.Index].PlayerStates[fascistPlayers[i].Index].NameStatus = "fascist"
				game.SeatedPlayers[hitler.Index].PlayerStates[fascistPlayers[i].Index].NotificationStatus = "fascist"
			}
		}

		SendInProgressGameUpdate(game)
	})

	time.AfterFunc(5000*time.Millisecond, func() {
		for i := range game.SeatedPlayers {
			game.SeatedPlayers[i].PlayerStates[i].CardStatus.Flipped = false

			for j := range game.SeatedPlayers[i].PlayerStates {
				game.SeatedPlayers[i].PlayerStates[j].NotificationStatus = ""
			}
		}

		SendInProgressGameUpdate(game)
	})

	time.AfterFunc(5200*time.Millisecond, func() {
		for i := range game.PublicPlayerStates {
			game.PublicPlayerStates[i].CardStatus.CardDisplayed = false
		}

		SendInProgressGameUpdate(game)
	})

	time.AfterFunc(5400*time.Millisecond, func() {
		for i := range game.SeatedPlayers {
			for j := range game.SeatedPlayers[i].PlayerStates {
				game.SeatedPlayers[i].PlayerStates[j].CardStatus = types.CardStatus{
					CardBack: "",
				}
			}
		}

		StartElection(game, -1)
	})
}

func StartGame(game *types.GamePrivate) {
	game.GamePublic.GameState.TracksFlipped = true
	game.GamePublic.GeneralGameSettings.LivingPlayerCount = game.GamePublic.PlayerCount

	for i := 0; i < game.PlayerCount-1; i++ {
		j := utils.RandInt(uint32(i), uint32(game.PlayerCount))
		player := game.GamePublic.PublicPlayerStates[i]
		game.GamePublic.PublicPlayerStates[i] = game.GamePublic.PublicPlayerStates[j]
		game.GamePublic.PublicPlayerStates[j] = player
		game.GamePublic.PlayerMap[game.GamePublic.PublicPlayerStates[i].ID] = int(i)
		game.GamePublic.PlayerMap[game.GamePublic.PublicPlayerStates[j].ID] = int(j)
		game.GamePublic.PublicPlayerStates[i].Index = i
		game.GamePublic.PublicPlayerStates[j].Index = int(j)
	}

	for i := range game.GamePublic.PublicPlayerStates {
		for j := range game.GamePublic.PublicPlayerStates[i].PlayerStates {
			game.GamePublic.PublicPlayerStates[i].PlayerStates[j].CardStatus.CardBack = ""
		}
	}

	game.SeatedPlayers = make([]types.PlayerState, len(game.GamePublic.PublicPlayerStates))
	copy(game.SeatedPlayers, game.GamePublic.PublicPlayerStates)
	fmt.Println("Players", game.SeatedPlayers[0].CardStatus.CardBack == "", game.GamePublic.PublicPlayerStates[0].CardStatus.CardBack == "")
	Countdown(game, 5)
}
